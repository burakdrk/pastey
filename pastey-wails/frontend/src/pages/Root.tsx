import Sidebar from "@/components/sidebar";
import { useAtom } from "jotai";
import { globalState } from "@/lib/store";
import { useMemo } from "react";
import Clipboard from "@/pages/Clipboard";

function Root() {
  const [pageStack] = useAtom(globalState.pageStack);
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
