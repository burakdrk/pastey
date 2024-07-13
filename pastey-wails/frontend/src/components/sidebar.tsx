import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from "@/components/ui/tooltip";
import { Button } from "@/components/ui/button";
import { Home, CircleUserRound, ClipboardList, LifeBuoy, Settings } from "lucide-react";
import { BrowserOpenURL } from "../../wailsjs/runtime";
import { globalState } from "@/lib/store";
import { useAtom } from "jotai";
import { pages } from "@/lib/types";

function Sidebar() {
  const [pageStack, setPageStack] = useAtom(globalState.pageStack);

  function changeToPage(page: pages) {
    setPageStack([page]);
  }

  function highlightSelected(page: pages) {
    if (pageStack[pageStack.length - 1] === page) {
      return " bg-muted";
    }

    return "";
  }

  return (
    <aside className="z-20 flex h-full flex-col border-r">
      <nav className="grid gap-1 p-2">
        <TooltipProvider>
          <Tooltip>
            <TooltipTrigger asChild>
              <Button
                variant="ghost"
                size="icon"
                className={"rounded-lg" + highlightSelected("home")}
                aria-label="Home"
                onClick={() => changeToPage("home")}
              >
                <Home className="size-5" />
              </Button>
            </TooltipTrigger>
            <TooltipContent side="right" sideOffset={5}>
              Home
            </TooltipContent>
          </Tooltip>
          <Tooltip>
            <TooltipTrigger asChild>
              <Button
                variant="ghost"
                size="icon"
                className={"rounded-lg" + highlightSelected("clipboard")}
                aria-label="Clipboard"
                onClick={() => changeToPage("clipboard")}
              >
                <ClipboardList className="size-5" />
              </Button>
            </TooltipTrigger>
            <TooltipContent side="right" sideOffset={5}>
              Clipboard
            </TooltipContent>
          </Tooltip>
          <Tooltip>
            <TooltipTrigger asChild>
              <Button
                variant="ghost"
                size="icon"
                className={"rounded-lg" + highlightSelected("account")}
                aria-label="Account"
                onClick={() => changeToPage("account")}
              >
                <CircleUserRound className="size-5" />
              </Button>
            </TooltipTrigger>
            <TooltipContent side="right" sideOffset={5}>
              Account
            </TooltipContent>
          </Tooltip>
        </TooltipProvider>
      </nav>
      <nav className="mt-auto grid gap-1 p-2">
        <TooltipProvider>
          <Tooltip>
            <TooltipTrigger asChild>
              <Button
                variant="ghost"
                size="icon"
                className="mt-auto rounded-lg"
                aria-label="Help"
                onClick={() => BrowserOpenURL("https://github.com/burakdrk/pastey/issues")}
              >
                <LifeBuoy className="size-5" />
              </Button>
            </TooltipTrigger>
            <TooltipContent side="right" sideOffset={5}>
              Help
            </TooltipContent>
          </Tooltip>
          <Tooltip>
            <TooltipTrigger asChild>
              <Button
                variant="ghost"
                size="icon"
                className={"mt-auto rounded-lg" + highlightSelected("settings")}
                aria-label="Settings"
                onClick={() => changeToPage("settings")}
              >
                <Settings className="size-5" />
              </Button>
            </TooltipTrigger>
            <TooltipContent side="right" sideOffset={5}>
              Settings
            </TooltipContent>
          </Tooltip>
        </TooltipProvider>
      </nav>
    </aside>
  );
}

export default Sidebar;
