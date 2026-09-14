import type { NextConfig } from "next";

import { envSchema, validateEnv } from "./config/env";

validateEnv(envSchema);

const nextConfig: NextConfig = {
  /* config options here */
  reactCompiler: true,
};

export default nextConfig;
