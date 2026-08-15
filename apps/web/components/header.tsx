import Link from "next/link";

import { Button } from "@loreline/ui/components/button";

import { Logo } from "./logo";

export const Header = () => {
  return (
    <header className="relative z-30">
      <div className="mx-auto flex h-18 max-w-304 items-center justify-between px-5 sm:px-8">
        <div className="flex items-center gap-7">
          <Logo />
          <nav
            aria-label="Homepage sections"
            className="hidden items-center gap-1 text-sm md:flex"
          >
            <a href="#how" className="rounded-full px-4 py-2 hover:bg-card">
              How it works
            </a>
            <a href="#voice" className="rounded-full px-4 py-2 hover:bg-card">
              Voice
            </a>
            <a
              href="#sideboard"
              className="rounded-full px-4 py-2 hover:bg-card"
            >
              Sideboard
            </a>
          </nav>
        </div>
        <div className="flex items-center gap-1.5">
          <Button
            variant="ghost"
            nativeButton={false}
            render={<Link href="/sign-in" prefetch />}
          >
            Sign in
          </Button>
          <Button
            nativeButton={false}
            render={<Link href="/library" prefetch />}
          >
            Open a book
          </Button>
        </div>
      </div>
    </header>
  );
};
