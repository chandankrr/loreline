import { z } from "zod";

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
});

export const ZRefreshResponse = z.object({
  accessToken: z.string(),
});
