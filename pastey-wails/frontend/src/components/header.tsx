import { Badge } from "./ui/badge";
import { WindowToggleMaximise, WindowMinimise, WindowHide, EventsOn, EventsOff } from "../../wailsjs/runtime";
import { useEffect, useState } from "react";
import { GetConnectionStatus, ConnectToWS } from "../../wailsjs/go/backend/App";

function Header() {
  const [isConnected, setIsConnected] = useState(false);
  const [isConnecting, setIsConnecting] = useState(true);

  let paddingString = "pl-20";

  if (window.navigator.userAgent.includes("Windows")) {
    paddingString = "pl-4";
  }

  useEffect(() => {
    EventsOn("ws:disconnected", () => {
      setIsConnected(false);
    });

    GetConnectionStatus().then((res) => {
      setIsConnected(res);
      setIsConnecting(false);
    });

    return () => {
      EventsOff("ws:disconnected");
    };
  }, []);

  async function connect() {
    if (isConnecting || isConnected) {
      return;
    }

    setIsConnecting(true);

    const err = await ConnectToWS();

    if (!err.error) {
      setIsConnected(true);
    }

    setIsConnecting(false);
  }

  return (
    <header
      // @ts-ignore
      style={{ "--wails-draggable": "drag" }}
      className={"sticky top-0 z-10 flex h-[42px] items-center gap-1 border-b bg-white dark:bg-black " + paddingString}
      onDoubleClick={() => WindowToggleMaximise()}
    >
      <h1 className="text-xl font-semibold flex-1">Pastey</h1>
      <Badge
        variant="outline"
        className={`${isConnecting ? "bg-yellow-500" : isConnected ? "bg-green-700" : "bg-red-600"} mr-3 ${
          !isConnecting && !isConnected && "cursor-pointer"
        }`}
        // @ts-ignore
        style={{ "--wails-draggable": "no-drag" }}
        onDoubleClick={(e) => e.stopPropagation()}
        onClick={() => connect()}
      >
        {isConnecting ? "Connecting" : isConnected ? "Connected" : "Disconnected"}
      </Badge>
      {window.navigator.userAgent.includes("Windows") && (
        <div
          className="flex h-full"
          // @ts-ignore
          style={{ "--wails-draggable": "no-drag" }}
          onDoubleClick={(e) => e.stopPropagation()}
        >
          <div
            className="w-[48px] h-full flex justify-center items-center bg-transparent hover:bg-white hover:bg-opacity-10 active:bg-opacity-20"
            onClick={() => WindowMinimise()}
          >
            <svg width="10px" height="10px" x="0px" y="0px" viewBox="0 0 10.2 1" className="dark:fill-white">
              <rect x="0" y="50%" width="10.2" height="1" />
            </svg>
          </div>
          <div
            className="w-[48px] h-full flex justify-center items-center bg-transparent hover:bg-white hover:bg-opacity-10 active:bg-opacity-20"
            onClick={() => WindowToggleMaximise()}
          >
            <svg width="10px" height="10px" viewBox="0 0 10 10" className="dark:fill-white">
              <path d="M0,0v10h10V0H0z M9,9H1V1h8V9z" />
            </svg>
          </div>
          <div
            className="w-[48px] h-full flex justify-center items-center bg-transparent hover:bg-red-600 active:bg-red-900"
            onClick={() => WindowHide()}
          >
            <svg width="10px" height="10px" viewBox="0 0 10 10" className="dark:fill-white">
              <polygon points="10.2,0.7 9.5,0 5.1,4.4 0.7,0 0,0.7 4.4,5.1 0,9.5 0.7,10.2 5.1,5.8 9.5,10.2 10.2,9.5 5.8,5.1" />
            </svg>
          </div>
        </div>
      )}
    </header>
  );
}

export default Header;
