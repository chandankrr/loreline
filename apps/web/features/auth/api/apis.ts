import type { TApiClient } from "@/api/api-client";

import { authStore } from "../stores/auth-store";
import type {
  TCurrentUserResponse,
  TForgotPasswordPayload,
  TForgotPasswordResponse,
  TLoginPayload,
  TLoginResponse,
  TLogoutResponse,
  TRefreshResponse,
  TRegisterPayload,
  TRegisterResponse,
  TResendEmailVerificationPayload,
  TResendEmailVerificationResponse,
  TResetPasswordPayload,
  TResetPasswordResponse,
  TVerifyEmailPayload,
  TVerifyEmailResponse,
} from "./types";

export const fetchSession = async ({
  api,
}: {
  api: TApiClient;
}): Promise<{ accessToken: string; user: TCurrentUserResponse } | null> => {
  try {
    const { accessToken } = await refresh({ api });
    authStore.setAccessToken(accessToken);
    const user = await currentUser({ api });
    return { accessToken, user };
  } catch {
    authStore.clearAuth();
    return null;
  }
};

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

export const login = async ({
  api,
  data,
}: {
  api: TApiClient;
  data: TLoginPayload;
}): Promise<TLoginResponse> => {
  const res = await api.Authentication.login({ body: data });

  if (res.status === 200) {
    return res.body;
  }

  throw res.body;
};

export const logout = async ({
  api,
}: {
  api: TApiClient;
}): Promise<TLogoutResponse> => {
  const res = await api.Authentication.logout({
    body: undefined,
  });

  if (res.status === 204) {
    return res.body;
  }

  throw res.body;
};

export const refresh = async ({
  api,
}: {
  api: TApiClient;
}): Promise<TRefreshResponse> => {
  const res = await api.Authentication.refreshToken({
    body: undefined,
  });

  if (res.status === 200) {
    return res.body;
  }

  throw res.body;
};

export const currentUser = async ({
  api,
}: {
  api: TApiClient;
}): Promise<TCurrentUserResponse> => {
  const res = await api.Authentication.me();

  if (res.status === 200) {
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
