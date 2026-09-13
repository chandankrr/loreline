import type { ServerInferRequest } from "@ts-rest/core";

import type { apiContract } from "@loreline/openapi/contracts";

export type TRequests = ServerInferRequest<typeof apiContract>;
