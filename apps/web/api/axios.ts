import axios, { type AxiosError, type AxiosRequestConfig } from "axios";

import type { RefreshResponse } from "@/features/auth/api";
import { authStore } from "@/features/auth/stores/auth-store";

declare module "axios" {
  export interface AxiosRequestConfig {
    _retry?: boolean;
  }
}

export const axiosInstance = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL,
  headers: { "Content-Type": "application/json" },
});

axiosInstance.interceptors.request.use((config) => {
  const token = authStore.getAccessToken();

  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }

  return config;
});

// Single-flight refresh so concurrent 401s share one refresh call
let refreshPromise: Promise<string | null> | null = null;

const refreshAccessToken = async (): Promise<string | null> => {
  if (!refreshPromise) {
    refreshPromise = axios
      .post<RefreshResponse>("/api/auth/refresh")
      .then((res) => {
        const token = res.data.accessToken ?? null;
        if (token) {
          authStore.setAccessToken(token);
        } else {
          authStore.clearAuth();
        }
        return token;
      })
      .catch(() => {
        authStore.clearAuth();
        return null;
      })
      .finally(() => {
        refreshPromise = null;
      });
  }

  return refreshPromise;
};

axiosInstance.interceptors.response.use(
  (response) => response,
  async (error: AxiosError) => {
    const originalRequest = error.config as AxiosRequestConfig | undefined;

    if (
      error.response?.status === 401 &&
      originalRequest &&
      !originalRequest._retry
    ) {
      originalRequest._retry = true;
      const newToken = await refreshAccessToken();

      if (newToken) {
        originalRequest.headers = {
          ...originalRequest.headers,
          Authorization: `Bearer ${newToken}`,
        };

        return axiosInstance(originalRequest);
      }

      if (typeof window !== "undefined") window.location.href = "/sign-in";
    }

    return Promise.reject(error);
  },
);
