"use client";

import { type PropsWithChildren, useEffect } from "react";

import { useApiClient } from "@/api/api-client";

import { fetchSession } from "../api";
import { useAuthStore } from "../stores/auth-store";

export const AuthProvider = ({ children }: PropsWithChildren) => {
  const api = useApiClient();
  const setAuth = useAuthStore((s) => s.setAuth);
  const clearAuth = useAuthStore((s) => s.clearAuth);
  const setInitialized = useAuthStore((s) => s.setInitialized);

  useEffect(() => {
    let cancelled = false;

    const initializeAuth = async () => {
      const session = await fetchSession({ api });

      if (cancelled) return;

      if (session) {
        setAuth(session.accessToken, session.user);
      } else {
        // Don't wipe out a session established by another auth flow
        // while this initial session check was running.
        const currentUser = useAuthStore.getState().user;

        if (!currentUser) {
          clearAuth();
        }
      }

      setInitialized(true);
    };

    initializeAuth();

    return () => {
      cancelled = true;
    };
  }, [api, setAuth, clearAuth, setInitialized]);

  return <>{children}</>;
};
