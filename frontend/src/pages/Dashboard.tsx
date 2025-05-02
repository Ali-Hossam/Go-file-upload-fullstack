import BackgroundLeft from "../components/BackgroundLeft";
import BackgroundRight from "../components/BackgroundRight";
import StudentsTable from "../components/StudentsTable";
import { useEffect, useState } from "react";
import { useSearchParams } from "react-router";
import { Record, SearchParamsKeys } from "../types/shared";
import NameSearch from "../components/NameSearch";
import Options from "../components/Options";

const DEFAULT_PAGE = 1;
const DEFAULT_PAGE_SIZE = 100;

export default function Dashboard() {
  const [searchParams] = useSearchParams();
  const page = searchParams.get(SearchParamsKeys.PAGE) || DEFAULT_PAGE;
  const pageSize =
    searchParams.get(SearchParamsKeys.PAGE_SIZE) || DEFAULT_PAGE_SIZE;
  const sortBy = searchParams.get(SearchParamsKeys.SORT_BY) || "";
  const sortOrder = searchParams.get(SearchParamsKeys.SORT_ORDER) || "";
  const name = searchParams.get(SearchParamsKeys.NAME) || "";
  const subject = searchParams.get(SearchParamsKeys.SUBJECT) || "";

  const [records, setRecords] = useState<Record[]>([]);
  const [count, setCount] = useState<number>(0);

  const serverUrl = import.meta.env.VITE_SERVER_URL;

  useEffect(() => {
    const capitalizedSortBy =
      sortBy.toLowerCase().charAt(0).toUpperCase() + sortBy.slice(1);
    const sortOrderSmallCase = sortOrder.toLowerCase();
    const capitalizedSubject =
      subject.charAt(0).toUpperCase() + subject.slice(1);

    fetch(
      `http://${serverUrl}/api/students?page=${page}&size=${pageSize}&sort_by=${capitalizedSortBy}&sort_order=${sortOrderSmallCase}&name=${name}&subject=${capitalizedSubject}`,
    )
      .then((response) => response.json())
      .then((data) => {
        setRecords(data.records);
        setCount(data.count);
      });
  }, [page, pageSize, sortBy, sortOrder, name, subject]);

  return (
    <main className="p-10 h-screen relative overflow-hidden selection:bg-teal-800 selection:text-teal-50">
      <div className="w-full h-full gap-2.5 items-center justify-center flex">
        <div className="h- h-full flex flex-col gap-2.5">
          <div className=" flex flex-col gap-12 justify-center bg-teal-50/20 backdrop-blur-md border-teal-600 px-20 py-8 border rounded-4xl">
            <NameSearch />
          </div>
          <div className="w-3xl overflow-auto h-full flex flex-col gap-12 justify-center bg-teal-50/20 backdrop-blur-md border-teal-600 p-2 border rounded-4xl">
            <StudentsTable records={records} />
          </div>
        </div>

        <div className="h-full flex flex-col gap-12 bg-teal-50/20 backdrop-blur-md border-teal-600 px-10 py-8  border rounded-4xl">
          <Options
            numOfPages={Math.ceil(count / parseInt(pageSize as string))}
          />
        </div>
      </div>
      <BackgroundLeft className="absolute saturate-0 opacity-10 -z-10 -left-12 bottom-0 rotate-12" />
      <BackgroundRight className="absolute -z-10 saturate-0  top-0 -right-36 rotate-12 opacity-16" />
    </main>
  );
}
