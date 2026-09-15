import axios, {
  type AxiosError,
  type AxiosInstance,
  type AxiosRequestConfig,
  type InternalAxiosRequestConfig,
} from "axios";

import type { TRefreshResponse } from "@/features/auth/api";
import { authStore } from "@/features/auth/stores/auth-store";

import { API_ENDPOINTS, apiConfig } from "./constants";

interface ApiRequestConfig extends AxiosRequestConfig {
  _retry?: boolean;
}

const refreshClient = axios.create(apiConfig);

let refreshRequest: Promise<string | null> | null = null;

const refreshAccessToken = async (): Promise<string | null> => {
  if (!refreshRequest) {
    refreshRequest = refreshClient
      .post<TRefreshResponse>(API_ENDPOINTS.REFRESH)
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
        refreshRequest = null;
      });
  }

  return refreshRequest;
};

export const attachAccessToken = (config: InternalAxiosRequestConfig) => {
  const token = authStore.getAccessToken();

  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }

  return config;
};

export const resolveUnauthorizedRequest = async (
  error: AxiosError,
  apiClient: AxiosInstance,
) => {
  const originalRequest = error.config as ApiRequestConfig | undefined;

  const isAuthRequest = originalRequest?.url?.includes("/auth");

  if (!originalRequest || originalRequest._retry || isAuthRequest) {
    return Promise.reject(error);
  }

  originalRequest._retry = true;

  const newToken = await refreshAccessToken();

  if (newToken) {
    originalRequest.headers = {
      ...originalRequest.headers,
      Authorization: `Bearer ${newToken}`,
    };

    return apiClient(originalRequest);
  }

  if (typeof window !== "undefined") {
    window.location.href = "/sign-in";
  }

  return Promise.reject(error);
};
