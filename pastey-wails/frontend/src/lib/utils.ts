import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";
import { parseISO } from "date-fns";
import { format } from "date-fns-tz";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function formatDateToLocalTime(isoDate: string): string {
  const parsedDate = parseISO(isoDate);

  return format(parsedDate, "hh:mm a · MMM d, yyyy", {
    timeZone: Intl.DateTimeFormat().resolvedOptions().timeZone,
  });
}

export const sleep = (ms: number) => new Promise((r) => setTimeout(r, ms));
