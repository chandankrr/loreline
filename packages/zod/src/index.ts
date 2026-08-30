import { extendZodWithOpenApi } from "@anatine/zod-openapi";
import { z } from "zod";

extendZodWithOpenApi(z);

export * from "./auth/index.js";
export * from "./health.js";
export * from "./user/index.js";
export * from "./utils.js";
