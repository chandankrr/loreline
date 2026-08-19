import { extendZodWithOpenApi } from "@anatine/zod-openapi";
import { generateOpenApi } from "@ts-rest/open-api";
import { z } from "zod";

import { apiContract } from "./contracts/index.js";

extendZodWithOpenApi(z);

type SecurityRequirementObject = {
  [key: string]: string[];
};

export type OperationMapper = NonNullable<
  Parameters<typeof generateOpenApi>[2]
>["operationMapper"];

const hasSecurity = (
  metadata: unknown,
): metadata is { openApiSecurity: SecurityRequirementObject[] } => {
  return (
    !!metadata && typeof metadata === "object" && "openApiSecurity" in metadata
  );
};

const operationMapper: OperationMapper = (operation, appRoute) => ({
  ...operation,
  ...(hasSecurity(appRoute.metadata)
    ? {
        security: appRoute.metadata.openApiSecurity,
      }
    : {}),
});

export const OpenAPI = Object.assign(
  generateOpenApi(
    apiContract,
    {
      openapi: "3.0.2",
      info: {
        version: "1.0.0",
        title: "Loreline REST API - Documentation",
        description: "Documentation for the Loreline REST API",
      },
      servers: [
        {
          url: "http://localhost:8080",
          description: "Local Server",
        },
      ],
    },
    {
      operationMapper,
      setOperationId: true,
    },
  ),
  {
    components: {
      securitySchemes: {
        bearerAuth: {
          type: "http",
          scheme: "bearer",
          bearerFormat: "JWT",
        },
        "x-service-token": {
          type: "apiKey",
          name: "x-service-token",
          in: "header",
        },
      },
    },
  },
);
