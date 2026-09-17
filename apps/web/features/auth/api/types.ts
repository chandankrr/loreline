import type { ServerInferResponseBody } from "@ts-rest/core";

import type { apiContract } from "@loreline/openapi/contracts";

import type { TRequests } from "@/api/types";

export type TRegisterPayload = TRequests["Authentication"]["register"]["body"];
export type TRegisterResponse = ServerInferResponseBody<
  typeof apiContract.Authentication.register,
  201
>;

export type TLoginPayload = TRequests["Authentication"]["login"]["body"];
export type TLoginResponse = ServerInferResponseBody<
  typeof apiContract.Authentication.login,
  200
>;

export type TLogoutPayload = TRequests["Authentication"]["logout"]["body"];
export type TLogoutResponse = ServerInferResponseBody<
  typeof apiContract.Authentication.logout,
  204
>;

export type TRefreshPayload =
  TRequests["Authentication"]["refreshToken"]["body"];
export type TRefreshResponse = ServerInferResponseBody<
  typeof apiContract.Authentication.refreshToken,
  200
>;

export type TVerifyEmailPayload =
  TRequests["Authentication"]["verifyEmail"]["body"];
export type TVerifyEmailResponse = ServerInferResponseBody<
  typeof apiContract.Authentication.verifyEmail,
  200
>;

export type TResendEmailVerificationPayload =
  TRequests["Authentication"]["resendEmailVerification"]["body"];
export type TResendEmailVerificationResponse = ServerInferResponseBody<
  typeof apiContract.Authentication.resendEmailVerification,
  200
>;

export type TForgotPasswordPayload =
  TRequests["Authentication"]["forgotPassword"]["body"];
export type TForgotPasswordResponse = ServerInferResponseBody<
  typeof apiContract.Authentication.forgotPassword,
  200
>;

export type TResetPasswordPayload =
  TRequests["Authentication"]["resetPassword"]["body"];
export type TResetPasswordResponse = ServerInferResponseBody<
  typeof apiContract.Authentication.resetPassword,
  200
>;

export type TCurrentUserResponse = ServerInferResponseBody<
  typeof apiContract.Authentication.me
>;
