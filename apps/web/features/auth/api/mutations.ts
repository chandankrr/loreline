import { useMutation } from "@tanstack/react-query";

import { useApiClient } from "@/api/api-client";
import { getApiErrorCode, showApiErrorToast } from "@/api/utils";

import { authStore } from "../stores/auth-store";
import {
  currentUser,
  forgotPassword,
  login,
  logout,
  register,
  resendEmailVerification,
  resetPassword,
  verifyEmail,
} from "./apis";
import type {
  TForgotPasswordPayload,
  TLoginPayload,
  TRegisterPayload,
  TResendEmailVerificationPayload,
  TResetPasswordPayload,
  TVerifyEmailPayload,
} from "./types";

export const useRegister = () => {
  const api = useApiClient();

  return useMutation({
    mutationFn: ({ body }: { body: TRegisterPayload }) =>
      register({ api, data: body }),
    onError: (err) => {
      if (getApiErrorCode(err) === "EMAIL_ALREADY_IN_USE") return;
      showApiErrorToast(err, "Failed to register");
    },
  });
};

export const useLogin = () => {
  const api = useApiClient();

  return useMutation({
    mutationFn: async ({ body }: { body: TLoginPayload }) => {
      const { accessToken } = await login({ api, data: body });
      authStore.setAccessToken(accessToken);
      const user = await currentUser({ api });
      return { accessToken, user };
    },
    onSuccess: ({ accessToken, user }) => {
      authStore.setAuth(accessToken, user);
    },
    onError: (err) => {
      showApiErrorToast(err, "Failed to log in");
    },
  });
};

export const useLogout = () => {
  const api = useApiClient();

  return useMutation({
    mutationFn: () => logout({ api }),
    onSettled: () => {
      authStore.clearAuth();
    },
  });
};

export const useVerifyEmail = () => {
  const api = useApiClient();

  return useMutation({
    mutationFn: ({ body }: { body: TVerifyEmailPayload }) =>
      verifyEmail({ api, data: body }),
    onError: (err) => {
      const code = getApiErrorCode(err);
      if (code === "INVALID_CODE" || code === "CODE_EXPIRED") return;
      showApiErrorToast(err, "Failed to verify email");
    },
  });
};

export const useResendEmailVerification = () => {
  const api = useApiClient();

  return useMutation({
    mutationFn: ({ body }: { body: TResendEmailVerificationPayload }) =>
      resendEmailVerification({ api, data: body }),
    onError: (err) => {
      showApiErrorToast(err, "Failed to resend verification email");
    },
  });
};

export const useForgotPassword = () => {
  const api = useApiClient();

  return useMutation({
    mutationFn: ({ body }: { body: TForgotPasswordPayload }) =>
      forgotPassword({ api, data: body }),
    onError: (err) => {
      showApiErrorToast(err, "Failed to send reset password link");
    },
  });
};

export const useResetPassword = () => {
  const api = useApiClient();

  return useMutation({
    mutationFn: ({ body }: { body: TResetPasswordPayload }) =>
      resetPassword({ api, data: body }),
    onError: (err) => {
      showApiErrorToast(err, "Failed to reset password");
    },
  });
};
