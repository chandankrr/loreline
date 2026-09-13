import { type NextRequest, NextResponse } from "next/server";

import { REFRESH_TOKEN_COOKIE_NAME } from "@/lib/constants";

export async function POST(request: NextRequest) {
  try {
    const refreshToken = request.cookies.get(REFRESH_TOKEN_COOKIE_NAME)?.value;
    const authHeader = request.headers.get("authorization");

    // Logout request to the original backend to revoke refresh token
    const backendResponse = await fetch(
      `${process.env.NEXT_PUBLIC_API_URL}/api/v1/auth/logout`,
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Cookie: `${REFRESH_TOKEN_COOKIE_NAME}=${refreshToken}`,
          ...(authHeader ? { Authorization: authHeader } : {}),
        },
      },
    );

    if (!backendResponse.ok) {
      const data = await backendResponse.json().catch(() => null);
      return NextResponse.json(data, { status: backendResponse.status });
    }

    const response = new NextResponse(null, {
      status: backendResponse.status,
    });

    // Clear refresh token cookie
    response.cookies.delete(REFRESH_TOKEN_COOKIE_NAME);

    return response;
  } catch (error) {
    console.error("Error logging out:", error);

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
