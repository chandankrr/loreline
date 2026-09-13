import { type NextRequest, NextResponse } from "next/server";

import { REFRESH_TOKEN_COOKIE_NAME } from "@/lib/constants";

import type { TLoginResponse } from "@/features/auth/api";

export async function POST(request: NextRequest) {
  try {
    const body = await request.json();

    // Login request to the original backend
    const backendResponse = await fetch(
      `${process.env.NEXT_PUBLIC_API_URL}/api/v1/auth/login`,
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(body),
      },
    );

    const data: TLoginResponse = await backendResponse.json();

    if (!backendResponse.ok) {
      return NextResponse.json(data, { status: backendResponse.status });
    }

    const { refreshToken, ...safeData } = data;

    const response = NextResponse.json(safeData, {
      status: backendResponse.status,
    });

    // Store refresh token in HttpOnly cookies
    if (refreshToken) {
      response.cookies.set(REFRESH_TOKEN_COOKIE_NAME, refreshToken, {
        httpOnly: true,
        secure: process.env.NODE_ENV === "production",
        sameSite: "strict",
        path: "/",
        maxAge: 60 * 60 * 24 * 7, // 7 days
      });
    }

    return response;
  } catch (error) {
    console.error("Error logging in:", error);

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
