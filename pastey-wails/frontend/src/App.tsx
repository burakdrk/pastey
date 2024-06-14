import { useState } from "react";
import { Greet } from "../wailsjs/go/main/App";
import { Button } from "./components/ui/button";

function App() {
  const [resultText, setResultText] = useState("Please enter your name below 👇");
  const [name, setName] = useState("");
  const updateName = (e: any) => setName(e.target.value);
  const [count, setCount] = useState(0);

  function greet() {
    Greet(name).then((result: string) => setResultText(result));
  }

  return (
    <div className="min-h-screen bg-white grid place-items-center mx-auto py-8">
      <div className="text-2xl font-bold flex flex-col items-center space-y-4">
        <h1>Vite + React + TS + Tailwind + shadcn/ui</h1>
        <Button onClick={() => setCount(count + 1)}>Count up ({count})</Button>
      </div>
    </div>
  );
}

export default App;
