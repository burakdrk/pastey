import { Button } from "@/components/ui/button";
import { Table, TableHeader, TableRow, TableHead, TableBody, TableCell } from "@/components/ui/table";
import { useEffect, useState } from "react";
import { DeleteEntry, GetEntries } from "../../wailsjs/go/backend/App";
import { models } from "wailsjs/go/models";
import { formatDateToLocalTime } from "@/lib/utils";
import { EventsOn, EventsOff } from "../../wailsjs/runtime/runtime";
import { RefreshCcw } from "lucide-react";
import { useAtom } from "jotai";
import { globalState } from "@/lib/store";
import { EntryEvent } from "@/lib/types";

function Clipboard() {
  async function getEntries() {
    setLoading(true);
    const res = await GetEntries();

    if (res.error.error) {
      console.error(res.error.error);
      return;
    }

    setEntries(res.entries);
    setLoading(false);
  }

  async function deleteEntry(id: string) {
    setLoading(true);
    const res = await DeleteEntry(id);

    if (res.error) {
      console.error(res.error);
      return;
    }

    setEntries(entries.filter((entry) => entry.entry_id !== id));
    setLoading(false);
  }

  async function deleteAllEntries() {
    setLoading(true);

    for (const entry of entries) {
      const res = await DeleteEntry(entry.entry_id);

      if (res.error) {
        console.error(res.error);
        return;
      }
    }

    setEntries([]);
    setLoading(false);
  }

  const [entries, setEntries] = useState<models.Entry[]>([]);
  const [loading, setLoading] = useState(true);
  const [deviceId] = useAtom(globalState.deviceId);

  useEffect(() => {
    getEntries();

    EventsOn("ws:entry", (data: EntryEvent) => {
      console.log(data);
    });

    return () => {
      EventsOff("ws:entry");
    };
  }, []);

  return (
    <main className="flex flex-1 flex-col gap-4 p-4 lg:gap-6 lg:p-6">
      <div className="flex items-center">
        <h1 className="text-lg font-semibold md:text-2xl flex-1">Clipboard</h1>
        <Button onClick={() => getEntries()} disabled={loading} variant={"outline"}>
          <RefreshCcw size={20} />
        </Button>
        <Button className="ml-4" variant={"destructive"} onClick={() => deleteAllEntries()} disabled={loading}>
          Delete All
        </Button>
      </div>

      <Table>
        <TableHeader>
          <TableRow>
            <TableHead className="w-[200px]">From</TableHead>
            <TableHead>Content</TableHead>
            <TableHead className="w-[100px]"></TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {entries.map((entry, index) => (
            <TableRow key={index}>
              <TableCell>
                <div className="font-medium">{entry.from_device_name}</div>
                <div className="hidden text-sm text-muted-foreground md:inline">
                  {formatDateToLocalTime(entry.created_at)}
                  {entry.from_device_id === deviceId && <div>This device</div>}
                </div>
              </TableCell>
              <TableCell className="break-all">{entry.encrypted_data}</TableCell>
              <TableCell className="flex gap-2 justify-end flex-col">
                <Button variant={"outline"} onClick={() => navigator.clipboard.writeText(entry.encrypted_data)}>
                  Copy
                </Button>
                <Button variant={"destructive"} onClick={() => deleteEntry(entry.entry_id)} disabled={loading}>
                  Delete
                </Button>
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </main>
  );
}

export default Clipboard;
