import { z } from "zod";

export const ZUser = z.object({
  id: z.string().uuid(),
  email: z.string().email(),
  name: z.string(),
  emailVerified: z.boolean(),
  image: z.string().nullable(),
  createdAt: z.string(),
  updatedAt: z.string(),
});

export const ZAuthUser = ZUser.omit({
  createdAt: true,
  updatedAt: true,
});
