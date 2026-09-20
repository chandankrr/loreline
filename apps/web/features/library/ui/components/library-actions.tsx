import { FolderPlusIcon, PlusIcon } from "lucide-react";

import { Button } from "@loreline/ui/components/button";

export const LibraryActions = () => {
  return (
    <div className="flex items-center gap-2">
      <Button size="lg" variant="outline" className="h-10">
        <FolderPlusIcon />
        New Stack
      </Button>
      <Button size="lg" className="h-10">
        <PlusIcon />
        Add Book
      </Button>
    </div>
  );
};
