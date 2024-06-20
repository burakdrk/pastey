import { Badge } from "./ui/badge";
import { WindowToggleMaximise } from "../../wailsjs/runtime";

function Header() {
  return (
    <header
      // @ts-ignore
      style={{ "--wails-draggable": "drag" }}
      className="sticky top-0 z-10 flex h-[42px] items-center gap-1 border-b bg-white dark:bg-black pr-4 pl-20"
      onDoubleClick={() => WindowToggleMaximise()}
    >
      <h1 className="text-xl font-semibold flex-1">pastey</h1>
      <Badge variant="outline" className="bg-green-700">
        Connected
      </Badge>
    </header>
  );
}

export default Header;
