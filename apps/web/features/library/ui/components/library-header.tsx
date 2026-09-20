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

import { Logo } from "@/components/logo";

import { useLogout } from "@/features/auth/api";
import { useAuth } from "@/features/auth/hooks/use-auth";

export const LibraryHeader = () => {
  const { user } = useAuth();
  const { mutate, isPending } = useLogout();

  const initials = user?.name
    ?.split(" ")
    .map((name) => name[0])
    .join("")
    .slice(0, 2)
    .toUpperCase();

  return (
    <header className="sticky top-0 z-40 border-b bg-background/88 backdrop-blur-xl">
      <div className="mx-auto flex h-17 max-w-7xl items-center justify-between px-4">
        <Logo />

        {user ? (
          <DropdownMenu>
            <DropdownMenuTrigger
              render={
                <Button variant="ghost" className="h-10 px-1.5 sm:pr-2.5" />
              }
            >
              <Avatar className="size-7">
                {user.image ? (
                  <AvatarImage
                    src={user.image}
                    alt={`${user.name}'s profile picture`}
                    referrerPolicy="no-referrer"
                  />
                ) : null}
                <AvatarFallback className="bg-primary text-[0.65rem] text-primary-foreground">
                  {initials}
                </AvatarFallback>
              </Avatar>
              <span className="hidden max-w-32 truncate text-sm sm:block">
                {user.name.trim().split(/\s+/)[0]}
              </span>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end" className="w-56">
              <div className="px-2 py-1.5">
                <p className="font-medium text-sm">{user.name}</p>
                <p className="truncate text-muted-foreground text-xs">
                  {user.email}
                </p>
              </div>
              <DropdownMenuSeparator />
              <DropdownMenuItem disabled={isPending} onClick={() => mutate()}>
                <LogOutIcon />
                Sign out
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        ) : null}
      </div>
    </header>
  );
};
