import { initContract } from "@ts-rest/core";

import { authenticationContract } from "./auth.js";
import { healthContract } from "./health.js";

const c = initContract();

export const apiContract = c.router({
  Health: healthContract,
  Authentication: authenticationContract,
});
