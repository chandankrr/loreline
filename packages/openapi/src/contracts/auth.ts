import { initContract } from "@ts-rest/core";
import { z } from "zod";

import {
  ZLoginPayload,
  ZLoginResponse,
  ZRefreshResponse,
  ZRegisterPayload,
  ZUser,
} from "@loreline/zod";

import { getSecurityMetadata } from "../utils.js";

const c = initContract();

const metadata = getSecurityMetadata();

export const authenticationContract = c.router(
  {
    register: {
      summary: "Register",
      path: "/register",
      method: "POST",
      description: "Creates a new user account",
      body: ZRegisterPayload.pick({
        email: true,
        name: true,
        password: true,
      }),
      responses: {
        201: ZUser,
      },
    },

    login: {
      summary: "Login",
      path: "/login",
      method: "POST",
      description: "Authenticates a user",
      body: ZLoginPayload.pick({
        email: true,
        password: true,
      }),
      responses: {
        200: ZLoginResponse,
      },
    },

    logout: {
      summary: "Logout",
      path: "/logout",
      method: "POST",
      description: "Logs out the authenticated user",
      body: z.void(),
      responses: {
        204: z.void(),
      },
      metadata: metadata,
    },

    refresh: {
      summary: "Refresh",
      path: "/refresh",
      method: "POST",
      description: "Refreshes the user's authentication tokens",
      body: z.void(),
      responses: {
        200: ZRefreshResponse,
      },
    },
  },
  {
    pathPrefix: "/api/v1/auth",
  },
);
