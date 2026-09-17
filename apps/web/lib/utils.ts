export function getSafeRedirect(
  redirect: string | null,
  fallback = "/library",
) {
  if (!redirect) return fallback;
  if (!redirect.startsWith("/") || redirect.startsWith("//")) {
    return fallback;
  }
  return redirect;
}
