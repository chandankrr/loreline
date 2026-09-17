import type { PropsWithChildren } from "react";

import { Toaster } from "@loreline/ui/components/toast";
import { TooltipProvider } from "@loreline/ui/components/tooltip";

import { AuthProvider } from "@/features/auth/providers/auth-provider";

import { QueryProvider } from "./query-provider";

export const AppProvider = ({ children }: PropsWithChildren) => {
  return (
    <QueryProvider>
      <AuthProvider>
        <TooltipProvider>{children}</TooltipProvider>
        <Toaster />
      </AuthProvider>
    </QueryProvider>
  );
};
