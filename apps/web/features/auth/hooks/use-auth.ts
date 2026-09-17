import { useAuthStore } from "../stores/auth-store";

export const useAuth = () => {
  const user = useAuthStore((s) => s.user);
  const isInitialized = useAuthStore((s) => s.isInitialized);

  return { user, isAuthenticated: !!user, isInitialized };
};
