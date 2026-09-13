import { useMutation } from "@tanstack/react-query";

import { useApiClient } from "@/api/api-client";
import { getApiErrorCode, showApiErrorToast } from "@/api/utils";

import {
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
  return useMutation({
    mutationFn: ({ body }: { body: TLoginPayload }) => login({ data: body }),
    onError: (err) => {
      showApiErrorToast(err, "Failed to login");
    },
  });
};

export const useLogout = () => {
  return useMutation({
    mutationFn: logout,
    onError: (err) => {
      showApiErrorToast(err, "Failed to logout");
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
