import { atom } from "jotai";

export const globalState = {
  isLoggedIn: atom(false),
  pageStack: atom<pages[]>(["clipboard"]),
};

export type pages = "home" | "settings" | "clipboard" | "account";
