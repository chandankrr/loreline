import { type NextRequest, NextResponse } from "next/server";

import { REFRESH_TOKEN_COOKIE_NAME } from "@/lib/constants";

import type { TRefreshResponse } from "@/features/auth/api";

export async function POST(request: NextRequest) {
  try {
    const refreshToken = request.cookies.get(REFRESH_TOKEN_COOKIE_NAME)?.value;

    if (!refreshToken) {
      const response = NextResponse.json(
        {
          code: "UNAUTHORIZED",
          message: "Missing refresh token",
          status: 401,
          override: false,
          errors: null,
          action: null,
        },
        { status: 401 },
      );

      response.cookies.delete(REFRESH_TOKEN_COOKIE_NAME);

      return response;
    }

    // Refresh request to the original backend
    // to issue new access + refresh tokens (token rotation)
    const backendResponse = await fetch(
      `${process.env.NEXT_PUBLIC_API_URL}/api/v1/auth/refresh`,
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Cookie: `${REFRESH_TOKEN_COOKIE_NAME}=${refreshToken}`,
        },
      },
    );

    const data: TRefreshResponse = await backendResponse.json();

    if (!backendResponse.ok) {
      const response = NextResponse.json(data, {
        status: backendResponse.status,
      });

      // Invalidate cookie if refresh token is expired or revoked
      response.cookies.delete(REFRESH_TOKEN_COOKIE_NAME);

      return response;
    }

    const { refreshToken: newRefreshToken, ...safeData } = data;

    const response = NextResponse.json(safeData, {
      status: backendResponse.status,
    });

    // Set newly rotated refreshToken in HttpOnly cookie
    if (newRefreshToken) {
      response.cookies.set(REFRESH_TOKEN_COOKIE_NAME, newRefreshToken, {
        httpOnly: true,
        secure: process.env.NODE_ENV === "production",
        sameSite: "strict",
        path: "/",
        maxAge: 60 * 60 * 24 * 7, // 7 days
      });
    }

    return response;
  } catch (error) {
    console.error("Error refreshing token:", error);

    return NextResponse.json(
      {
        code: "INTERNAL_SERVER_ERROR",
        message: "Internal server error",
        status: 500,
        override: false,
        errors: null,
        action: null,
      },
      { status: 500 },
    );
  }
}
