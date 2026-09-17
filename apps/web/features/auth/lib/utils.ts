import { API_ENDPOINTS } from "@/api/constants";

export const getOAuthUrl = (provider: "google" | string = "google") => {
  const baseURL = process.env.NEXT_PUBLIC_API_URL;
  return `${baseURL}${API_ENDPOINTS.OAUTH(provider)}`;
};
