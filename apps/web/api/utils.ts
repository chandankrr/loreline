import { toast } from "@loreline/ui/components/toast";

export const getApiErrorCode = (error: unknown): string | undefined =>
  error &&
  typeof error === "object" &&
  "code" in error &&
  typeof error.code === "string"
    ? error.code
    : undefined;

export const showApiErrorToast = (error: unknown, fallbackMessage: string) => {
  const message =
    error &&
    typeof error === "object" &&
    "message" in error &&
    typeof error.message === "string"
      ? error.message
      : fallbackMessage;

  console.error(message, error);

  toast.add({
    title: message,
    type: "error",
  });
};
