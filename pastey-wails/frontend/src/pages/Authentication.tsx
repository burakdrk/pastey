import LoginForm from "@/components/forms/login-form";
import { ModeToggle } from "@/components/mode-toggle";

function Authentication() {
  return (
    <>
      <div className="absolute top-0 right-0 px-4 py-4 mt-9">
        <ModeToggle />
      </div>
      <div className="flex items-center justify-center py-12 h-full">
        <div className="mx-auto grid w-[350px] gap-6">
          <div className="grid gap-2 text-center">
            <h1 className="text-3xl font-bold mb-4">Login</h1>
            <p className="text-balance text-muted-foreground">Enter your email below to login to your account</p>
          </div>

          <LoginForm />

          <div className="mt-4 text-center text-sm">
            Don&apos;t have an account?{" "}
            <a href="#" className="underline">
              Sign up
            </a>
          </div>
        </div>
      </div>
    </>
  );
}

export default Authentication;
