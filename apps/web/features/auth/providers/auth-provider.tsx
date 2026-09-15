"use client";

import { type PropsWithChildren, useEffect, useRef } from "react";

import { useApiClient } from "@/api/api-client";

import { fetchSession } from "../api";
import { useAuthStore } from "../stores/auth-store";

export const AuthProvider = ({ children }: PropsWithChildren) => {
  const api = useApiClient();
  const setAuth = useAuthStore((s) => s.setAuth);
  const clearAuth = useAuthStore((s) => s.clearAuth);
  const setInitialized = useAuthStore((s) => s.setInitialized);

  const sessionPromiseRef = useRef<ReturnType<typeof fetchSession> | null>(
    null,
  );

  useEffect(() => {
    let cancelled = false;

    if (!sessionPromiseRef.current) {
      sessionPromiseRef.current = fetchSession({ api });
    }

    sessionPromiseRef.current.then((session) => {
      if (cancelled) return;
      if (session) {
        setAuth(session.accessToken, session.user);
      } else {
        clearAuth();
      }
      setInitialized(true);
    });

    return () => {
      cancelled = true;
    };
  }, [api, clearAuth, setAuth, setInitialized]);

  return <>{children}</>;
};
