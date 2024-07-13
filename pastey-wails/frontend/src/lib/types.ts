export type pages = "home" | "settings" | "clipboard" | "account";

export type EntryEvent = {
  encrypted_data: string;
  from_device_id: number;
  to_device_id: number;
  user_id: number;
};
