import BackgroundLeft from "../components/BackgroundLeft";
import BackgroundRight from "../components/BackgroundRight";
import FileUpload from "../components/FileUpload";
import Logo from "../components/Logo";

export default function Home() {
  return (
    <main className="p-4 sm:p-6 md:p-8 lg:p-10 min-h-screen relative overflow-hidden selection:bg-teal-800 selection:text-teal-50">
      <div className="flex items-center min-h-[calc(100vh-2rem)] sm:min-h-[calc(100vh-3rem)] justify-center">
        <div className="w-full max-w-[95%] sm:max-w-[90%] md:max-w-[85%] lg:max-w-3xl flex flex-col gap-6 sm:gap-8 md:gap-12 justify-center bg-teal-50/20 backdrop-blur-xl border-teal-600 px-4 sm:px-8 md:px-12 lg:px-20 py-8 sm:py-12 md:py-16 lg:py-20 border rounded-2xl sm:rounded-3xl md:rounded-4xl">
          <Logo />

          <div className="flex w-full items-center justify-center">
            <FileUpload />
          </div>
        </div>
      </div>

      <BackgroundLeft className="absolute saturate-0 hidden sm:block opacity-10 -z-10 -left-12 bottom-0 rotate-12" />
      <BackgroundRight className="absolute -z-10 saturate-0 top-0 -right-36 rotate-12 opacity-16" />
    </main>
  );
}
