import { useMemo } from "react";

import { initClient } from "@ts-rest/core";
import { isAxiosError } from "axios";

import { apiContract } from "@loreline/openapi/contracts";

import { axiosInstance } from "./axios";

type Headers = Awaited<
  ReturnType<NonNullable<Parameters<typeof initClient>[1]["api"]>>
>["headers"];

export const useApiClient = ({ isBlob = false }: { isBlob?: boolean } = {}) =>
  useMemo(
    () =>
      initClient(apiContract, {
        baseUrl: "",
        baseHeaders: {
          "Content-Type": "application/json",
        },
        api: async ({ path, method, headers, body }) => {
          try {
            const result = await axiosInstance.request({
              method,
              url: path,
              headers,
              data: body,
              ...(isBlob && {
                responseType: "blob",
              }),
            });

            return {
              status: result.status,
              body: result.data,
              headers: result.headers as unknown as Headers,
            };
          } catch (error) {
            if (isAxiosError(error)) {
              const response = error.response;
              return {
                status: response?.status ?? 500,
                body: response?.data ?? { message: "Internal server error" },
                headers: (response?.headers as unknown as Headers) ?? {},
              };
            }

            throw error;
          }
        },
      }),
    [isBlob],
  );

export type TApiClient = ReturnType<typeof useApiClient>;
