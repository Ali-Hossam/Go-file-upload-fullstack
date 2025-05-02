export default function Logo() {
  return (
    <div className="flex flex-col gap-4 w-full items-center justify-center">
      <h1 className="font-nova font-extrabold text-6xl text-teal-600">
        <span className="text-teal-950">GO</span>Sheet
      </h1>
      <p className="text-black/40 max-w-2xl text-center">
        Effortlessly upload, visualize, and manage your CSV data with GoSheet —
        the fast, reliable Go-powered platform that transforms your spreadsheets
        into actionable insights stored securely in your database
      </p>
    </div>
  );
}
