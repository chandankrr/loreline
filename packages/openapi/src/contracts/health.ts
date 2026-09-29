import { initContract } from "@ts-rest/core";

import { ZHealthResponse } from "@loreline/zod";

const c = initContract();

export const healthContract = c.router({
  getHealth: {
    summary: "Get Health",
    path: "/healthz",
    method: "GET",
    description: "Get health status",
    responses: {
      200: ZHealthResponse,
    },
  },
});
