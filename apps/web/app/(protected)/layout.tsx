"use client";

import { useRouter } from "next/navigation";
import { useEffect } from "react";

import { useAuth } from "@/features/auth/hooks/use-auth";

export default function ProtectedLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const { user, isInitialized } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (isInitialized && !user) router.replace("/sign-in");
  }, [isInitialized, user, router]);

  // if (!isInitialized) return <PageSkeleton />;
  if (!user) return null;

  return <>{children}</>;
}
