import { z } from "zod";

import { ZAuthUser } from "../user/index.js";

export const ZOAuthProvider = z.enum(["google"]);

export const ZRegisterPayload = z.object({
  email: z.string().email(),
  name: z.string(),
  password: z.string(),
});

export const ZLoginPayload = z.object({
  email: z.string().email(),
  password: z.string(),
});

export const ZLoginResponse = z.object({
  accessToken: z.string(),
  refreshToken: z.string(),
  user: ZAuthUser,
});

export const ZRefreshResponse = z.object({
  accessToken: z.string(),
  refreshToken: z.string(),
  user: ZAuthUser,
});

export const ZVerifyEmailPayload = z.object({
  email: z.string().email(),
  code: z.string(),
});

export const ZResendVerificationCodePayload = z.object({
  email: z.string().email(),
});

export const ZForgotPasswordPayload = z.object({
  email: z.string().email(),
});

export const ZResetPasswordPayload = z.object({
  token: z.string(),
  password: z.string(),
});
