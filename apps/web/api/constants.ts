import type { CreateAxiosDefaults } from "axios";

export const API_ENDPOINTS = {
  REFRESH: "/api/v1/auth/refresh",
} as const;

export const apiConfig = {
  baseURL: process.env.NEXT_PUBLIC_API_URL,
  withCredentials: true,
} satisfies CreateAxiosDefaults;
