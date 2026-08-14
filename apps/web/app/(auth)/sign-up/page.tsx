import { Logo } from "@/components/logo";

import { SignUpForm } from "@/features/auth/ui/components/sign-up-form";

export default function SignUpPage() {
  return (
    <div className="flex min-h-[calc(100vh-1.5rem)] flex-col bg-background p-5 sm:p-8">
      <div className="flex items-center justify-between lg:justify-end">
        <Logo className="lg:hidden" />
      </div>
      <SignUpForm />
    </div>
  );
}
