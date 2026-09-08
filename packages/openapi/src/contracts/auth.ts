import { initContract } from "@ts-rest/core";
import { z } from "zod";

import {
  ZForgotPasswordPayload,
  ZLoginPayload,
  ZLoginResponse,
  ZMessageResponse,
  ZOAuthProvider,
  ZRefreshResponse,
  ZRegisterPayload,
  ZResendVerificationCodePayload,
  ZResetPasswordPayload,
  ZUser,
  ZVerifyEmailPayload,
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

    refreshToken: {
      summary: "Refresh Token",
      path: "/refresh",
      method: "POST",
      description:
        "Refreshes the user's authentication tokens using a refresh token stored in an HTTP-only cookie + token rotation",
      body: z.void(),
      responses: {
        200: ZRefreshResponse,
      },
    },

    verifyEmail: {
      summary: "Email Verification",
      path: "/email/verify",
      method: "POST",
      description:
        "Verifies the user's email address using the verification code",
      body: ZVerifyEmailPayload,
      responses: {
        200: ZMessageResponse,
      },
    },

    resendEmailVerification: {
      summary: "Resend Email Verification",
      path: "/email/resend",
      method: "POST",
      description:
        "Resends a new email verification code to the user's registered email address",
      body: ZResendVerificationCodePayload,
      responses: {
        200: ZMessageResponse,
      },
    },

    forgotPassword: {
      summary: "Forgot Password",
      path: "/password/forgot",
      method: "POST",
      description:
        "Initiates password reset process by sending a reset password email to the user's registered email address",
      body: ZForgotPasswordPayload,
      responses: {
        200: ZMessageResponse,
      },
    },

    resetPassword: {
      summary: "Reset Password",
      path: "/password/reset",
      method: "POST",
      description: "Resets the user's password",
      body: ZResetPasswordPayload,
      responses: {
        200: ZMessageResponse,
      },
    },

    oauthBegin: {
      summary: "Begin OAuth",
      path: "/:provider",
      method: "GET",
      description:
        "Starts the OAuth authentication flow for the specified provider",
      pathParams: z.object({
        provider: ZOAuthProvider,
      }),
      responses: {
        307: z.void(),
      },
    },

    oauthCallback: {
      summary: "OAuth Callback",
      path: "/:provider/callback",
      method: "GET",
      description:
        "Completes the OAuth authentication flow for the specified provider and redirects to the frontend",
      pathParams: z.object({
        provider: ZOAuthProvider,
      }),
      responses: {
        307: z.void(),
      },
    },
  },
  {
    pathPrefix: "/api/v1/auth",
  },
);
