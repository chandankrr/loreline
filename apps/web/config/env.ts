import { z } from "zod";

export const envSchema = z.object({
  NODE_ENV: z.string().default("development"),
  NEXT_PUBLIC_APP_URL: z.url().default("http://localhost:3000"),
  NEXT_PUBLIC_API_URL: z.url().default("http://localhost:8080"),
});

export const validateEnv = (schema: z.ZodType) => {
  try {
    schema.parse(process.env);
  } catch (error) {
    if (error instanceof z.ZodError) {
      console.log(
        " \x1b[1m\x1b[31m⚠\x1b[0m Environment variables validation failed\n",
      );

      error.issues.forEach((err) => {
        console.error(
          `   - \x1b[31m${err.path.join(".")}\x1b[0m: ${err.message}`,
        );
      });
    } else {
      console.error(
        " \x1b[31m⚠\x1b[0m Unexpected error during environment validation:",
        error,
      );
    }

    console.log();
    process.exit(1);
  }
};

declare global {
  namespace NodeJS {
    interface ProcessEnv extends z.infer<typeof envSchema> {}
  }
}
