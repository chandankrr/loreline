import { create } from "zustand";

import type { TCurrentUserResponse } from "../api";

type User = TCurrentUserResponse;

type AuthState = {
  accessToken: string | null;
  user: User | null;
  isInitialized: boolean; // has the initial refresh-on-load check finished?
  setAuth: (accessToken: string, user: User) => void;
  setAccessToken: (accessToken: string) => void;
  clearAuth: () => void;
  setInitialized: (v: boolean) => void;
};

export const useAuthStore = create<AuthState>((set) => ({
  accessToken: null,
  user: null,
  isInitialized: false,
  setAuth: (accessToken, user) => set({ accessToken, user }),
  setAccessToken: (accessToken) => set({ accessToken }),
  clearAuth: () => set({ accessToken: null, user: null }),
  setInitialized: (isInitialized) => set({ isInitialized }),
}));

export const authStore = {
  getAccessToken: () => useAuthStore.getState().accessToken,
  setAccessToken: (accessToken: string) =>
    useAuthStore.getState().setAccessToken(accessToken),
  clearAuth: () => useAuthStore.getState().clearAuth(),
};
