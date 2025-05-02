import { Record } from "../types/shared";

export default function StudentsTable({ records }: { records: Record[] }) {
  return (
    <div className="overflow-hidden h-full rounded-3xl border border-teal-600 backdrop-blur-sm bg-teal-50/30">
      <div className="overflow-y-auto h-full">
        <table className="w-full table-auto">
          <thead className="sticky top-0 bg-teal-700 text-white z-10">
            <tr>
              <th className="py-4 w-80 text-lg px-10 text-left font-medium">
                Name
              </th>
              <th className="py-4 text-lg px-10 text-left font-medium">
                Subject
              </th>
              <th className="py-4 text-lg px-10 text-left font-medium">
                Grade
              </th>
            </tr>
          </thead>
          <tbody>
            {records?.map((record, idx) => (
              <tr
                key={record.Student_id}
                className={`border-b border-teal-700/10 hover:bg-teal-700/10 transition-colors duration-300 ${idx % 2 === 0 ? "bg-neutral-200/20" : ""}`}
              >
                <td className="py-2 px-10 text-left">{record.Student_name}</td>
                <td className="py-2 px-10 text-left">{record.Subject}</td>
                <td className="py-2 px-10 text-left font-bold">
                  {record.Grade}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
