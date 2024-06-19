import { useState } from "react";
import { Greet } from "../wailsjs/go/main/App";
import { Button } from "./components/ui/button";
import { ThemeProvider } from "./components/theme-provider";
import { ModeToggle } from "./components/mode-toggle";
import { Authentication } from "./pages/Authentication";

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
      <Authentication />
      <div className="absolute top-0 right-0 px-4 py-4">
        <ModeToggle />
      </div>
    </ThemeProvider>
  );
}

export default App;
