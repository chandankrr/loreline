import { Logo } from "@/components/logo";

import { ResetPasswordForm } from "@/features/auth/ui/components/reset-password-form";

export default async function ResetPasswordPage({
  searchParams,
}: {
  searchParams: Promise<{ token?: string }>;
}) {
  const { token } = await searchParams;

  return (
    <div className="flex min-h-[calc(100vh-1.5rem)] flex-col bg-background p-5 sm:p-8">
      <div className="flex items-center justify-between lg:justify-end">
        <Logo className="lg:hidden" />
      </div>
      <ResetPasswordForm token={token ?? ""} />
    </div>
  );
}
