import { useSearchParams } from "react-router";

export default function NameSearch() {
  const [searchParams, setSearchParams] = useSearchParams();

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const newParams = new URLSearchParams(searchParams);
    newParams.set("name", e.target.value);
    setSearchParams(newParams);
  };

  return (
    <div className="flex flex-col items-center justify-center gap-4">
      <h1 className="text-2xl font-bold font-nova">Search by name</h1>
      <input
        onChange={handleInputChange}
        className="border-2 border-teal-800 transition-all text-lg rounded-2xl px-6 bg-white p-2 w-full focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-teal-800 focus-visible:ring-opacity-75"
      />
    </div>
  );
}
