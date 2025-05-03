import { ChevronLeft, ChevronRight } from "lucide-react";
import { useSearchParams } from "react-router";
import { SearchParamsKeys } from "../types/shared";
import { useState } from "react";

const subjects = {
  Mathematics: "Mathematics",
  Physics: "Physics",
  Chemistry: "Chemistry",
  Biology: "Biology",
  History: "History",
  EnglishLit: "English Literature",
  CompSci: "Computer Science",
  Art: "Art",
  Music: "Music",
  Geography: "Geography",
};

export default function Options({ numOfPages }: { numOfPages: number }) {
  const [searchParams, setSearchParams] = useSearchParams();
  const [page, setPage] = useState(1);

  const updateSearchParams = (param: string, value: string) => {
    const newParams = new URLSearchParams(searchParams);
    newParams.set(param, value);
    setSearchParams(newParams);
  };

  const handleSortBySelection = (
    event: React.ChangeEvent<HTMLSelectElement>,
  ) => {
    updateSearchParams(SearchParamsKeys.SORT_BY, event.target.value);
  };

  const handleSortOrderSelection = (
    event: React.ChangeEvent<HTMLInputElement>,
  ) => {
    updateSearchParams(SearchParamsKeys.SORT_ORDER, event.target.value);
  };

  const handleSubjectSelection = (
    event: React.ChangeEvent<HTMLSelectElement>,
  ) => {
    updateSearchParams(SearchParamsKeys.SUBJECT, event.target.value);
  };

  const handleItemsPerPageSelection = (
    event: React.ChangeEvent<HTMLSelectElement>,
  ) => {
    updateSearchParams(SearchParamsKeys.PAGE_SIZE, event.target.value);
  };

  const handlePrevPage = () => {
    const currentPage = parseInt(
      searchParams.get(SearchParamsKeys.PAGE) || "1",
      10,
    );
    if (currentPage > 1) {
      updateSearchParams(SearchParamsKeys.PAGE, (currentPage - 1).toString());
      setPage(currentPage - 1);
    }
  };

  const handleNextPage = () => {
    const currentPage = parseInt(
      searchParams.get(SearchParamsKeys.PAGE) || "1",
      10,
    );
    if (currentPage < numOfPages) {
      updateSearchParams(SearchParamsKeys.PAGE, (currentPage + 1).toString());
      setPage(currentPage + 1);
    }
  };

  const handlePageChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const page = parseInt(event.target.value, 10);
    if (page >= 1 && page <= numOfPages) {
      updateSearchParams(SearchParamsKeys.PAGE, page.toString());
      setPage(page);
    }
  };

  return (
    <div className="flex flex-col gap-4 sm:gap-6 md:gap-8 w-full">
      <h2 className="text-xl sm:text-2xl font-bold font-nova text-black">
        Table Options
      </h2>

      {/* Sorting Section */}
      <div className="space-y-2 sm:space-y-3 md:space-y-4">
        <h3 className="text-base sm:text-lg font-semibold text-teal-700">
          Sort By
        </h3>
        <div className="flex flex-col gap-2 sm:gap-3">
          <div className="w-full">
            <select
              className="w-full px-2 sm:px-3 md:px-4 py-1.5 sm:py-2 text-sm sm:text-base rounded-lg sm:rounded-xl border-2 border-teal-700 bg-white focus:outline-none focus:ring-2 focus:ring-teal-700 transition-all"
              onChange={handleSortBySelection}
            >
              <option value="">Select a column</option>
              <option value="student_name">Name</option>
              <option value="subject">Subject</option>
              <option value="grade">Grade</option>
            </select>
          </div>

          <div className="flex flex-wrap gap-2 sm:gap-4 items-center">
            <span className="text-xs sm:text-sm text-gray-700">
              Sort Order:
            </span>
            <div className="flex flex-wrap gap-2 sm:gap-3">
              <label className="flex items-center gap-1 sm:gap-2 text-xs sm:text-sm">
                <input
                  type="radio"
                  name="sortOrder"
                  value="asc"
                  className="text-teal-600 accent-teal-700"
                  onChange={handleSortOrderSelection}
                />
                <span>Ascending</span>
              </label>
              <label className="flex items-center gap-1 sm:gap-2 text-xs sm:text-sm">
                <input
                  type="radio"
                  name="sortOrder"
                  value="desc"
                  className="text-teal-600 accent-teal-700"
                  onChange={handleSortOrderSelection}
                />
                <span>Descending</span>
              </label>
            </div>
          </div>
        </div>
      </div>

      {/* Filter by Subject */}
      <div className="space-y-2 sm:space-y-3 md:space-y-4">
        <h3 className="text-base sm:text-lg font-semibold text-teal-700">
          Filter by Subject
        </h3>
        <div className="w-full">
          <select
            className="w-full bg-white px-2 sm:px-3 md:px-4 py-1.5 sm:py-2 text-sm sm:text-base rounded-lg sm:rounded-xl border-2 border-teal-700 transition-all focus:outline-none focus:ring-2 focus:ring-teal-700"
            onChange={handleSubjectSelection}
          >
            <option value="">All Subjects</option>
            {Object.entries(subjects).map(([key, value]) => (
              <option key={key} value={key}>
                {value}
              </option>
            ))}
          </select>
        </div>
      </div>

      {/* Pagination Section */}
      <div className="space-y-2 sm:space-y-3 md:space-y-4">
        <h3 className="text-base sm:text-lg font-semibold text-teal-700">
          Pagination
        </h3>
        <div className="flex flex-col gap-2 sm:gap-3">
          <div className="flex items-center flex-wrap gap-2 sm:gap-3">
            <span className="text-xs sm:text-sm text-gray-700">
              Items per page:
            </span>
            <select
              className="px-2 sm:px-3 bg-white py-1 text-xs sm:text-sm rounded-md sm:rounded-lg border-2 border-teal-700 transition-all focus:outline-none focus:ring-2 focus:ring-teal-700"
              onChange={handleItemsPerPageSelection}
              defaultValue={100}
            >
              <option value="50">50</option>
              <option value="100">100</option>
              <option value="250">250</option>
              <option value="500">500</option>
              <option value="1000">1000</option>
            </select>
          </div>

          <div className="flex items-center justify-between mt-1 sm:mt-2 bg-white rounded-lg sm:rounded-xl border-2 border-teal-700 p-1 sm:p-2">
            <button
              className="p-0.5 sm:p-1 rounded-full hover:bg-teal-700/70 hover:cursor-pointer hover:text-white transition-all"
              onClick={handlePrevPage}
            >
              <ChevronLeft size={16} className="sm:size-auto" />
            </button>

            <div className="flex items-center gap-1 sm:gap-2 text-xs sm:text-sm">
              <span>Page</span>
              <input
                min="1"
                value={page}
                className="w-10 sm:w-16 text-center rounded-md sm:rounded-lg border border-teal-700/70 py-0.5 sm:py-1 text-xs sm:text-sm focus:outline-none focus:ring-2 focus:ring-teal-700/70"
                onChange={handlePageChange}
              />
              <span>of {numOfPages}</span>
            </div>

            <button
              className="p-0.5 sm:p-1 rounded-full hover:bg-teal-700/70 hover:cursor-pointer hover:text-white transition-all"
              onClick={handleNextPage}
            >
              <ChevronRight size={16} className="sm:size-auto" />
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
