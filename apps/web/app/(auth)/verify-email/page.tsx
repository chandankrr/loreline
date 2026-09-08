import { Logo } from "@/components/logo";

import { VerifyEmailForm } from "@/features/auth/ui/components/verify-email-form";

export default async function VerifyEmailPage({
  searchParams,
}: {
  searchParams: Promise<{ email?: string }>;
}) {
  const { email } = await searchParams;

  return (
    <div className="flex min-h-[calc(100vh-1.5rem)] flex-col bg-background p-5 sm:p-8">
      <div className="flex items-center justify-between lg:justify-end">
        <Logo className="lg:hidden" />
      </div>
      <VerifyEmailForm email={email ?? ""} />
    </div>
  );
}
