import { Bird, CornerDownLeft, Mic, Paperclip, Rabbit, Turtle } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Textarea } from "@/components/ui/textarea";
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from "@/components/ui/tooltip";
import Header from "@/components/header";
import Sidebar from "@/components/sidebar";
import { useAtom } from "jotai";
import { globalState } from "@/lib/store";
import { useMemo } from "react";
import Clipboard from "@/pages/Clipboard";

function Root() {
  const [pageStack, setPageStack] = useAtom(globalState.pageStack);
  const peek = useMemo(() => pageStack[pageStack.length - 1], [pageStack]);

  function renderPage() {
    switch (peek) {
      case "home":
        return "home";
      case "settings":
        return "settings";
      case "clipboard":
        return <Clipboard />;
      case "account":
        return "account";
      default:
        return null;
    }
  }

  return (
    <>
      <div className="flex h-full w-full overflow-y-auto">
        <Sidebar />

        {renderPage()}
      </div>
    </>
  );
}

export default Root;
