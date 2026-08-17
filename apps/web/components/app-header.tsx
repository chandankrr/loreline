import Link from "next/link";

import { LogOutIcon } from "lucide-react";

import {
  Avatar,
  AvatarFallback,
  AvatarImage,
} from "@loreline/ui/components/avatar";
import { Button } from "@loreline/ui/components/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@loreline/ui/components/dropdown-menu";

import { Logo } from "./logo";
import { RealtimeModelBadge } from "./realtime-model-badge";

type AppHeaderProps = {
  variant?: "app" | "marketing";
};

export const AppHeader = ({ variant = "app" }: AppHeaderProps) => {
  const isMarketing = variant === "marketing";

  return (
    <header
      className={
        isMarketing
          ? "relative z-30 bg-background"
          : "sticky top-0 z-40 border-b bg-background/88 backdrop-blur-xl"
      }
    >
      <div
        className={
          isMarketing
            ? "mx-auto flex h-18 max-w-304 items-center justify-between px-5 sm:px-8"
            : "mx-auto flex h-17 max-w-7xl items-center justify-between px-4 sm:px-0"
        }
      >
        <div className="flex items-center gap-7">
          <Logo />
          {isMarketing && (
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
          )}
        </div>
        <div className="flex items-center gap-1.5">
          <RealtimeModelBadge className="mr-1 hidden sm:inline-flex" />
          {isMarketing ? (
            <>
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
            </>
          ) : (
            <DropdownMenu>
              <DropdownMenuTrigger
                render={
                  <Button
                    variant="ghost"
                    className="ml-1 h-10 gap-2 px-1.5 pr-2.5"
                  />
                }
              >
                <Avatar className="size-7">
                  <AvatarImage
                    src=""
                    alt="'s profile picture"
                    referrerPolicy="no-referrer"
                  />
                  <AvatarFallback className="bg-primary text-[0.65rem] text-primary-foreground">
                    C
                  </AvatarFallback>
                </Avatar>
                <span className="hidden max-w-32 truncate text-sm sm:block">
                  Chandan
                </span>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end" sideOffset={10} className="w-56">
                <div className="px-2 py-1.5">
                  <p className="font-medium text-sm">Chandan</p>
                  <p className="truncate text-muted-foreground text-xs">
                    chandankrr.91@gmail.com
                  </p>
                </div>
                <DropdownMenuSeparator />
                <DropdownMenuItem>
                  <LogOutIcon />
                  Sign out
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          )}
        </div>
      </div>
    </header>
  );
};
