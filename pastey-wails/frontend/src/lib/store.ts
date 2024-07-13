import { atom } from "jotai";
import { pages } from "./types";

export const globalState = {
  isLoggedIn: atom(false),
  pageStack: atom<pages[]>(["clipboard"]),
  deviceId: atom(0),
};
