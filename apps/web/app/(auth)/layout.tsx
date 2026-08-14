import { CheckIcon } from "lucide-react";

import { Logo } from "@/components/logo";

export default function AuthLayout({ children }: LayoutProps<"/">) {
  return (
    <div className="min-h-screen bg-background p-3 lg:grid lg:grid-cols-[1.05fr_0.95fr] lg:gap-3">
      <div className="relative hidden overflow-hidden rounded-[2.5rem] bg-card p-12 lg:flex lg:flex-col">
        <Logo />
        <div className="my-auto max-w-xl">
          <p className="mb-5 font-semibold text-brand-ink text-sm">
            Your private reading room
          </p>
          <h1 className="max-w-xl font-semibold text-6xl leading-[1.05] tracking-tighter">
            A better conversation starts with the page.
          </h1>
          <div className="mt-10 space-y-4 text-muted-foreground text-sm">
            {[
              "Your PDFs stay scoped to your account",
              "Page-first answers with optional book retrieval",
              "Voice, visual explanations, notes, and progress together",
            ].map((point) => (
              <div key={point} className="flex items-center gap-3">
                <span className="grid size-6 place-items-center rounded-full bg-brand-soft text-brand-ink">
                  <CheckIcon className="size-3.5" />
                </span>
                {point}
              </div>
            ))}
          </div>
        </div>
        <p className="text-muted-foreground text-xs">
          Loreline · Read deeply. Keep wondering.
        </p>
      </div>
      {children}
    </div>
  );
}
