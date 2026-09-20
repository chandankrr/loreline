import { useState } from "react";

import { ListFilterPlusIcon, SearchIcon } from "lucide-react";

import { Badge } from "@loreline/ui/components/badge";
import { Button } from "@loreline/ui/components/button";
import {
  DropdownMenu,
  DropdownMenuCheckboxItem,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuLabel,
  DropdownMenuRadioGroup,
  DropdownMenuRadioItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@loreline/ui/components/dropdown-menu";
import { Input } from "@loreline/ui/components/input";

const STATUS_OPTIONS = [
  { value: "reading", label: "Reading" },
  { value: "finished", label: "Finished" },
] as const;

const SORT_OPTIONS = [
  { value: "recent", label: "Recently added" },
  { value: "title", label: "Title (A–Z)" },
  { value: "author", label: "Author" },
] as const;

export const ShelfToolbar = () => {
  const [statusFilters, setStatusFilters] = useState<string[]>([]);
  const [sort, setSort] = useState<string>("recent");

  const toggleStatus = (value: string) => {
    setStatusFilters((prev) =>
      prev.includes(value) ? prev.filter((v) => v !== value) : [...prev, value],
    );
  };

  const activeFilterCount = statusFilters.length;

  return (
    <div className="flex items-center justify-between">
      <div className="flex min-w-0 flex-1 items-center gap-3">
        <div className="relative w-full max-w-sm">
          <SearchIcon className="pointer-events-none absolute top-1/2 left-3 size-4 -translate-y-1/2 text-muted-foreground" />
          <Input placeholder="Search your shelf" className="pl-9" />
        </div>
        <Badge variant="secondary" className="hidden sm:block">
          2 books
        </Badge>
      </div>
      <DropdownMenu>
        <DropdownMenuTrigger
          render={
            <Button variant="outline">
              Filters <ListFilterPlusIcon />
              {activeFilterCount > 0 && (
                <Badge variant="secondary" className="ml-1">
                  {activeFilterCount}
                </Badge>
              )}
            </Button>
          }
        />
        <DropdownMenuContent className="w-48">
          <DropdownMenuGroup>
            <DropdownMenuLabel>Status</DropdownMenuLabel>
            {STATUS_OPTIONS.map((option) => (
              <DropdownMenuCheckboxItem
                key={option.value}
                checked={statusFilters.includes(option.value)}
                onCheckedChange={() => toggleStatus(option.value)}
              >
                {option.label}
              </DropdownMenuCheckboxItem>
            ))}
          </DropdownMenuGroup>
          <DropdownMenuSeparator />
          <DropdownMenuRadioGroup value={sort} onValueChange={setSort}>
            <DropdownMenuLabel>Sort by</DropdownMenuLabel>
            {SORT_OPTIONS.map((option) => (
              <DropdownMenuRadioItem key={option.value} value={option.value}>
                {option.label}
              </DropdownMenuRadioItem>
            ))}
          </DropdownMenuRadioGroup>
        </DropdownMenuContent>
      </DropdownMenu>
    </div>
  );
};
