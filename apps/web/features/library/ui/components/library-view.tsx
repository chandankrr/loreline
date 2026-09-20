"use client";

import { LibraryActions } from "./library-actions";
import { LibraryHeader } from "./library-header";
import { ShelfSection } from "./shelf-section";
import { ShelfToolbar } from "./shelf-toolbar";

export const LibraryView = () => {
  return (
    <div className="min-h-screen bg-background">
      <LibraryHeader />

      <main className="mx-auto max-w-7xl px-4 py-10 sm:py-16">
        <div className="flex flex-col justify-between gap-7 sm:flex-row sm:items-end">
          <div>
            <p className="font-semibold text-brand-ink text-sm">
              Your reading room
            </p>
            <h1 className="mt-3 font-semibold text-5xl tracking-[-0.047em]">
              Library
            </h1>
            <p className="mt-3 text-muted-foreground">
              A shelf for books you want to understand, not merely finish.
            </p>
          </div>
          <LibraryActions />
        </div>

        <div className="mt-10 border-b pb-5">
          <ShelfToolbar />
        </div>

        <ShelfSection />
      </main>
    </div>
  );
};
