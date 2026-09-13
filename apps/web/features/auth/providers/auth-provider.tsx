"use client";

import { type PropsWithChildren, useEffect } from "react";

import { refreshSession } from "../api";
import { useAuthStore } from "../stores/auth-store";

export const AuthProvider = ({ children }: PropsWithChildren) => {
  const setAuth = useAuthStore((s) => s.setAuth);
  const clearAuth = useAuthStore((s) => s.clearAuth);
  const setInitialized = useAuthStore((s) => s.setInitialized);

  useEffect(() => {
    let cancelled = false;

    refreshSession().then((data) => {
      if (cancelled) return;
      if (data) {
        setAuth(data.accessToken, data.user);
      } else {
        clearAuth();
      }
      setInitialized(true);
    });

    return () => {
      cancelled = true;
    };
  }, [clearAuth, setAuth, setInitialized]);

  return <>{children}</>;
};
