import axios, { isAxiosError } from "axios";

import type { TApiClient } from "@/api/api-client";

import { authStore } from "../stores/auth-store";
import type {
  LoginResponse,
  RefreshResponse,
  TForgotPasswordPayload,
  TForgotPasswordResponse,
  TLoginPayload,
  TLogoutResponse,
  TRegisterPayload,
  TRegisterResponse,
  TResendEmailVerificationPayload,
  TResendEmailVerificationResponse,
  TResetPasswordPayload,
  TResetPasswordResponse,
  TVerifyEmailPayload,
  TVerifyEmailResponse,
} from "./types";

// ----------- Handled via Next.js routes -----------
let refreshPromise: Promise<RefreshResponse | null> | null = null;

export const refreshSession = (): Promise<RefreshResponse | null> => {
  if (!refreshPromise) {
    refreshPromise = axios
      .post<RefreshResponse>("/api/auth/refresh")
      .then((res) => res.data)
      .catch(() => null)
      .finally(() => {
        refreshPromise = null;
      });
  }

  return refreshPromise;
};

export const login = async ({
  data,
}: {
  data: TLoginPayload;
}): Promise<LoginResponse> => {
  try {
    const res = await axios.post<LoginResponse>("/api/auth/login", data);

    return res.data;
  } catch (error) {
    if (isAxiosError(error) && error.response?.data) {
      throw error.response.data;
    }

    throw error;
  }
};

export const logout = async (): Promise<TLogoutResponse> => {
  const token = authStore.getAccessToken();

  try {
    const res = await axios.post("/api/auth/logout", null, {
      headers: token ? { Authorization: `Bearer ${token}` } : undefined,
    });

    return res.data;
  } catch (error) {
    if (isAxiosError(error) && error.response?.data) {
      throw error.response.data;
    }

    throw error;
  }
};

// ----------- Handled via api client -----------

export const register = async ({
  api,
  data,
}: {
  api: TApiClient;
  data: TRegisterPayload;
}): Promise<TRegisterResponse> => {
  const res = await api.Authentication.register({ body: data });

  if (res.status === 201) {
    return res.body;
  }

  throw res.body;
};

export const verifyEmail = async ({
  api,
  data,
}: {
  api: TApiClient;
  data: TVerifyEmailPayload;
}): Promise<TVerifyEmailResponse> => {
  const res = await api.Authentication.verifyEmail({ body: data });

  if (res.status === 200) {
    return res.body;
  }

  throw res.body;
};

export const resendEmailVerification = async ({
  api,
  data,
}: {
  api: TApiClient;
  data: TResendEmailVerificationPayload;
}): Promise<TResendEmailVerificationResponse> => {
  const res = await api.Authentication.resendEmailVerification({
    body: data,
  });

  if (res.status === 200) {
    return res.body;
  }

  throw res.body;
};

export const forgotPassword = async ({
  api,
  data,
}: {
  api: TApiClient;
  data: TForgotPasswordPayload;
}): Promise<TForgotPasswordResponse> => {
  const res = await api.Authentication.forgotPassword({ body: data });

  if (res.status === 200) {
    return res.body;
  }

  throw res.body;
};

export const resetPassword = async ({
  api,
  data,
}: {
  api: TApiClient;
  data: TResetPasswordPayload;
}): Promise<TResetPasswordResponse> => {
  const res = await api.Authentication.resetPassword({ body: data });

  if (res.status === 200) {
    return res.body;
  }

  throw res.body;
};
