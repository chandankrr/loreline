import axios from "axios";

import { attachAccessToken, resolveUnauthorizedRequest } from "./axios-auth";
import { apiConfig } from "./constants";

export const axiosInstance = axios.create({
  ...apiConfig,
  headers: { "Content-Type": "application/json" },
});

axiosInstance.interceptors.request.use(attachAccessToken);

axiosInstance.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      return resolveUnauthorizedRequest(error, axiosInstance);
    }

    return Promise.reject(error);
  },
);
