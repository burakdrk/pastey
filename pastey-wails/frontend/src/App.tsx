import { useState } from "react";
import { Greet } from "../wailsjs/go/main/App";
import { Button } from "./components/ui/button";
import { ThemeProvider } from "./components/theme-provider";
import { ModeToggle } from "./components/mode-toggle";
import { Auth } from "./pages/Authentication";

function App() {
  const [resultText, setResultText] = useState("Please enter your name below 👇");
  const [name, setName] = useState("");
  const updateName = (e: any) => setName(e.target.value);
  const [count, setCount] = useState(0);

  function greet() {
    Greet(name).then((result: string) => setResultText(result));
  }

  return (
    <ThemeProvider defaultTheme="dark" storageKey="ui-theme">
      {/* <div className="min-h-screen grid place-items-center mx-auto py-8">
        <div className="text-2xl font-bold flex flex-col items-center space-y-4">
          <h1>Vite + React + TS + Tailwind + shadcn/ui</h1>
          <ModeToggle />
          <Button className="bg-primary" onClick={() => setCount(count + 1)}>
            Count up ({count})
          </Button>
        </div>
      </div> */}
      <Auth />
      <div className="absolute top-0 right-0 px-4 py-4">
        <ModeToggle />
      </div>
    </ThemeProvider>
  );
}

export default App;
