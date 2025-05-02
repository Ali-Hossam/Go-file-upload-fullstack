import BackgroundLeft from "../components/BackgroundLeft";
import BackgroundRight from "../components/BackgroundRight";
import FileUpload from "../components/FileUpload";
import Logo from "../components/Logo";

export default function Home() {
  return (
    <main className="p-10 h-screen relative overflow-hidden selection:bg-teal-800 selection:text-teal-50">
      <div className="flex items-center h-full justify-center">
        <div className="w-fit flex flex-col gap-12 justify-center bg-teal-50/20 backdrop-blur-xl border-teal-600 px-40 py-20 border rounded-4xl">
          <Logo />

          <div className="flex w-full items-center justify-center">
            <FileUpload />
          </div>
        </div>
      </div>

      <BackgroundLeft className="absolute saturate-0 opacity-10 -z-10 -left-12 bottom-0 rotate-12" />
      <BackgroundRight className="absolute -z-10 saturate-0  top-0 -right-36 rotate-12 opacity-16" />
    </main>
  );
}
