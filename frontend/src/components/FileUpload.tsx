import {
  CloudUpload,
  X,
  FileText,
  AlertCircle,
  CheckCircle2,
} from "lucide-react";
import { useRef, useState, ChangeEvent, DragEvent, useEffect } from "react";
import { Link } from "react-router";

interface ProcessStatus {
  Id: number;
  Percent: number;
  Timeleft: number;
  Error: string;
}

interface UploadStatus {
  inProgress: boolean;
  percent: number;
  timeLeft?: number;
  error: string | null;
  completed: boolean;
}

// [50% AI]
export default function FileUpload() {
  const [files, setFiles] = useState<File[]>([]);
  const [isDragging, setIsDragging] = useState(false);
  const fileInputRef = useRef<HTMLInputElement>(null);
  const [statusId, setStatusId] = useState<string | null>(null);
  const wsRef = useRef<WebSocket | null>(null);
  const [filesUploadStatus, setFilesUploadStatus] = useState<
    Map<number, ProcessStatus>
  >(new Map());
  const [uploadStatus, setUploadStatus] = useState<UploadStatus>({
    inProgress: false,
    percent: 0,
    timeLeft: 0,
    error: null,
    completed: false,
  });

  const uploadStatusRef = useRef<UploadStatus>(uploadStatus);

  useEffect(() => {
    uploadStatusRef.current = uploadStatus;
  }, [uploadStatus]);

  useEffect(() => {
    if (!statusId) return;
    if (uploadStatusRef.current.completed) return;

    setUploadStatus((prev) => ({
      ...prev,
      inProgress: true,
      percent: 0,
      completed: false,
    }));

    // Extract server URL without protocol
    const serverUrl = import.meta.env.VITE_SERVER_URL;

    // Establish WebSocket connection
    const ws = new WebSocket(`ws://${serverUrl}/api/upload/status/${statusId}`);

    ws.onopen = () => {
      console.log("WebSocket connection established");
    };

    ws.onmessage = (event) => {
      console.log("Received message:", event.data);

      function roundToOneDecimal(num: number) {
        return Math.round((num + Number.EPSILON) * 10) / 10;
      }

      try {
        const data: ProcessStatus = JSON.parse(event.data);
        data.Percent = roundToOneDecimal(data.Percent);
        data.Timeleft = Math.round(data.Timeleft);

        setFilesUploadStatus((prev) => {
          const newMap = new Map(prev).set(data.Id, data);

          if (newMap.size > 0) {
            // Calculate average progress and time
            let totalPercent = 0;
            let totalTimeLeft = 0;
            let allFilesComplete = true;

            newMap.forEach((status) => {
              totalPercent += status.Percent;
              totalTimeLeft += status.Timeleft;
              if (status.Percent < 100) allFilesComplete = false;
            });

            const avgPercent = Math.round(totalPercent / newMap.size);
            const avgTimeLeft = Math.round(totalTimeLeft / newMap.size);

            setUploadStatus((prev) => ({
              ...prev,
              percent: avgPercent,
              timeLeft: avgTimeLeft,
              inProgress: !allFilesComplete,
              completed: allFilesComplete,
            }));
          }

          return newMap;
        });
      } catch (err) {
        console.error("Error parsing WebSocket message:", err);
        setUploadStatus((prev) => ({
          ...prev,
          error: `Error parsing server message: ${err instanceof Error ? err.message : String(err)}`,
          inProgress: false,
        }));
      }
    };

    ws.onclose = (event) => {
      if (uploadStatusRef.current.completed) {
        return;
      }
      console.log(
        `WebSocket connection closed. Code: ${event.code}, Reason: '${event.reason}', WasClean: ${event.wasClean}`,
      );
      // Only show error if not completed successfully and not already showing error
      if (!uploadStatusRef.current.error) {
        setUploadStatus((prev) => ({
          ...prev,
          error: `Connection closed: ${event.reason || "Unknown reason"} (${event.code})`,
          inProgress: false,
        }));
      }
    };

    ws.onerror = (error) => {
      console.error("WebSocket error:", error);
      setUploadStatus((prev) => ({
        ...prev,
        error:
          "WebSocket connection error occurred. See browser console for details.",
        inProgress: false,
      }));
    };

    wsRef.current = ws;

    return () => {
      ws.close();
    };
  }, [statusId]);

  const handleDragEnter = (e: DragEvent<HTMLDivElement>) => {
    e.preventDefault();
    e.stopPropagation();
    setIsDragging(true);
  };

  const handleDragOver = (e: DragEvent<HTMLDivElement>) => {
    e.preventDefault();
    e.stopPropagation();
    setIsDragging(true);
  };

  const handleDragLeave = (e: DragEvent<HTMLDivElement>) => {
    e.preventDefault();
    e.stopPropagation();
    setIsDragging(false);
  };

  const handleDrop = (e: DragEvent<HTMLDivElement>) => {
    e.preventDefault();
    e.stopPropagation();
    setIsDragging(false);

    if (e.dataTransfer.files.length > 0) {
      const validFiles = Array.from(e.dataTransfer.files).filter((file) =>
        file.name.toLowerCase().endsWith(".csv"),
      );

      if (validFiles.length > 0) {
        setFiles((prev) => [...prev, ...validFiles]);
      }
    }
  };

  const handleFileChange = (e: ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files.length > 0) {
      const fileArray = Array.from(e.target.files);
      setFiles((prev) => [...prev, ...fileArray]);
    }
  };

  const handleClick = () => {
    fileInputRef.current?.click();
  };

  const removeFile = (index: number) => {
    const newFiles = [...files];
    newFiles.splice(index, 1);
    setFiles(newFiles);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (files.length === 0) {
      setUploadStatus({
        inProgress: false,
        percent: 0,
        error: "Please select at least one CSV file to upload",
        completed: false,
      });
      return;
    }

    // Validate file types on client side before sending
    const invalidFiles = files.filter(
      (file) => !file.name.toLowerCase().endsWith(".csv"),
    );
    if (invalidFiles.length > 0) {
      setUploadStatus({
        inProgress: false,
        percent: 0,
        error: `Invalid file type(s): ${invalidFiles.map((f) => f.name).join(", ")}. Only CSV files are allowed.`,
        completed: false,
      });
      return;
    }

    const formData = new FormData();

    for (let i = 0; i < files.length; i++) {
      formData.append("files", files[i]);
    }

    setUploadStatus({
      inProgress: true,
      percent: 0,
      error: null,
      completed: false,
    });

    try {
      // Extract server URL without protocol if it exists
      const serverUrl = import.meta.env.VITE_SERVER_URL;

      const apiUrl = `http://${serverUrl}/api/upload`;

      console.log("Sending upload to:", apiUrl);
      console.log(
        "Files being uploaded:",
        files.map((f) => ({ name: f.name, type: f.type, size: f.size })),
      );

      const response = await fetch(apiUrl, {
        method: "POST",
        body: formData,
      });

      if (response.ok) {
        const result = await response.json();
        console.log("Upload started with ID:", result.upload_id);
        setStatusId(result.upload_id);
      } else {
        const errorData = await response.text();
        console.error("Upload failed:", errorData);
        setUploadStatus({
          inProgress: false,
          percent: 0,
          error: `Upload failed: ${errorData || response.statusText}`,
          completed: false,
        });
      }
    } catch (error) {
      console.error("Error during upload:", error);
      setUploadStatus({
        inProgress: false,
        percent: 0,
        error: `Error during upload: ${error instanceof Error ? error.message : String(error)}`,
        completed: false,
      });
    }
  };

  // Format time remaining
  const formatTimeRemaining = (seconds: number): string => {
    if (!seconds || !isFinite(seconds) || seconds <= 0) return "calculating...";
    if (seconds < 60) return `${seconds} sec`;
    return `${Math.floor(seconds / 60)} min ${seconds % 60} sec`;
  };

  return (
    <form
      className="flex flex-col w-full max-w-full sm:max-w-lg md:max-w-xl lg:max-w-2xl gap-3 sm:gap-4 justify-center items-center"
      onSubmit={handleSubmit}
      method="POST"
    >
      <div className="w-full border border-gray-200 rounded-lg overflow-hidden bg-white/50">
        {/* Header section with upload area */}
        <div
          onClick={!uploadStatus.inProgress ? handleClick : undefined}
          onDragEnter={!uploadStatus.inProgress ? handleDragEnter : undefined}
          onDragOver={!uploadStatus.inProgress ? handleDragOver : undefined}
          onDragLeave={!uploadStatus.inProgress ? handleDragLeave : undefined}
          onDrop={!uploadStatus.inProgress ? handleDrop : undefined}
          className={`w-full flex items-center justify-center p-3 sm:p-4 md:p-6 border-b border-gray-200 transition-all
            ${isDragging ? "bg-teal-50 border-teal-300" : "hover:bg-gray-100"}
            ${uploadStatus.inProgress ? "opacity-75" : "cursor-pointer"}`}
        >
          <div className="flex items-center">
            <CloudUpload className="text-teal-600 mr-2 sm:mr-3 md:mr-4" size={24} strokeWidth={1.5} />
            <div>
              <h2 className="font-medium text-sm sm:text-base md:text-lg">
                {uploadStatus.inProgress
                  ? `Uploading ${files.length} file${files.length !== 1 ? "s" : ""}...`
                  : "Choose files or drag & drop"}
              </h2>
              <p className="text-xs sm:text-sm text-gray-500">
                {uploadStatus.inProgress
                  ? `${uploadStatus.percent}% complete`
                  : "CSV format only, up to 1GB"}
              </p>
            </div>
          </div>

          <input
            type="file"
            ref={fileInputRef}
            multiple
            accept=".csv"
            onChange={handleFileChange}
            className="hidden"
            disabled={uploadStatus.inProgress}
          />
        </div>

        {/* Progress bar for overall upload */}
        {uploadStatus.inProgress && (
          <div className="h-1 bg-gray-100 w-full">
            <div
              className="h-1 bg-teal-600 transition-all duration-300"
              style={{ width: `${uploadStatus.percent}%` }}
            ></div>
          </div>
        )}

        {/* File list section */}
        <div className="p-2 sm:p-3 md:p-4">
          {uploadStatus.error && (
            <div className="mb-3 sm:mb-4 p-2 sm:p-3 bg-red-50 text-red-800 rounded border border-red-200 flex items-start">
              <AlertCircle
                className="text-red-500 mr-2 flex-shrink-0 mt-0.5"
                size={16}
              />
              <div className="flex-1">
                <p className="text-xs sm:text-sm font-medium">{uploadStatus.error}</p>
                <button
                  type="button"
                  onClick={handleSubmit}
                  className="mt-2 bg-red-600 hover:bg-red-700 text-white font-medium py-1 px-2 sm:px-3 rounded text-xs sm:text-sm"
                  disabled={files.length === 0 || uploadStatus.inProgress}
                >
                  Retry Upload
                </button>
              </div>
            </div>
          )}

          {uploadStatus.completed && (
            <div className="mb-3 sm:mb-4 p-2 sm:p-3 bg-green-50 text-green-800 rounded border border-green-200 flex items-center">
              <CheckCircle2 className="text-green-500 mr-2" size={16} />
              <p className="text-xs sm:text-sm font-medium">
                All files uploaded successfully!
              </p>
            </div>
          )}

          {files.length > 0 ? (
            <div>
              <div className="flex justify-between items-center mb-2">
                <h3 className="font-medium text-xs sm:text-sm text-gray-700">
                  Files{" "}
                  {uploadStatus.inProgress &&
                    `(${uploadStatus.timeLeft ? formatTimeRemaining(uploadStatus.timeLeft) + " remaining" : "Processing..."})`}
                </h3>
                {!uploadStatus.inProgress && files.length > 0 && (
                  <button
                    type="submit"
                    className="text-[10px] sm:text-xs bg-teal-600 hover:bg-teal-700 text-white px-2 sm:px-3 py-1 rounded-md transition-colors"
                  >
                    Upload {files.length} file{files.length !== 1 ? "s" : ""}
                  </button>
                )}
              </div>

              <ul className="max-h-40 sm:max-h-48 md:max-h-60 overflow-y-auto divide-y divide-gray-100">
                {files.map((file, index) => {
                  // Find matching status for this file if available
                  const fileStatus = Array.from(
                    filesUploadStatus.entries(),
                  ).find(([id]) => id === index)?.[1];

                  return (
                    <li key={index} className="py-1 sm:py-2 flex items-center gap-2 sm:gap-3">
                      <FileText
                        size={16}
                        className="text-gray-400 flex-shrink-0"
                      />

                      <div className="flex-1 min-w-0">
                        <div className="flex justify-between items-center">
                          <span className="truncate font-medium text-xs sm:text-sm">
                            {file.name}
                          </span>
                          {!uploadStatus.inProgress && (
                            <button
                              type="button"
                              onClick={() => removeFile(index)}
                              className="text-red-500 hover:cursor-pointer p-1 hover:bg-red-50 rounded-full ml-2 flex-shrink-0"
                            >
                              <X size={12} />
                            </button>
                          )}
                        </div>

                        {/* Show progress bar if this file is being processed */}
                        {fileStatus && (
                          <div className="mt-1">
                            <div className="flex items-center justify-between text-[10px] sm:text-xs">
                              <div className="text-gray-500">
                                {fileStatus.Percent >= 100 ? (
                                  <span className="text-green-600 font-medium">
                                    Completed
                                  </span>
                                ) : fileStatus.Error ? (
                                  <span className="text-red-600">
                                    {fileStatus.Error}
                                  </span>
                                ) : (
                                  <span>
                                    {fileStatus.Percent}% •{" "}
                                    {formatTimeRemaining(fileStatus.Timeleft)}
                                  </span>
                                )}
                              </div>
                            </div>
                            <div className="w-full bg-gray-100 rounded-full h-1 mt-1">
                              <div
                                className={`h-1 rounded-full transition-all duration-300 ${
                                  fileStatus.Error
                                    ? "bg-red-500"
                                    : "bg-teal-600"
                                }`}
                                style={{ width: `${fileStatus.Percent}%` }}
                              ></div>
                            </div>
                          </div>
                        )}
                      </div>
                    </li>
                  );
                })}
              </ul>
            </div>
          ) : (
            <div className="text-center py-4 sm:py-6 md:py-8 text-gray-500 text-xs sm:text-sm">
              <p>No files selected</p>
            </div>
          )}
        </div>
      </div>

      {!uploadStatus.inProgress && !uploadStatus.completed && (
        <button
          type="submit"
          className={`rounded-md px-4 sm:px-6 md:px-8 font-medium border transition-colors text-sm sm:text-base md:text-lg py-1.5 sm:py-2 w-full
            ${
              files.length > 0
                ? "bg-teal-600 text-white hover:cursor-pointer hover:bg-teal-700 border-teal-700"
                : "hover:cursor-not-allowed text-gray-500 bg-gray-200 border-gray-300"
            }`}
          disabled={files.length === 0}
        >
          {files.length > 0
            ? `Upload ${files.length} File${files.length > 1 ? "s" : ""}`
            : "Select Files to Upload"}
        </button>
      )}
      {uploadStatus.completed && (
        <Link
          className="rounded-md px-4 sm:px-6 md:px-8 bg-teal-600 text-white text-center hover:cursor-pointer hover:bg-teal-700 border-teal-700 font-medium border transition-colors text-sm sm:text-base md:text-lg py-1.5 sm:py-2 w-full"
          to="/dashboard"
        >
          View Students Records &rarr;
        </Link>
      )}
    </form>
  );
}
