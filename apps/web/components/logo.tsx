import Link from "next/link";

import { cn } from "@loreline/ui/lib/utils";

type LogoProps = {
  compact?: boolean;
  className?: string;
};

export const Logo = ({ compact = false, className }: LogoProps) => {
  return (
    <Link
      href="/"
      aria-label="Loreline home"
      className={cn("inline-flex items-center", className)}
    >
      {compact ? (
        <LorelineMark />
      ) : (
        <svg
          aria-hidden="true"
          viewBox="0 0 139 42"
          className="h-[2.35rem] w-[7.85rem] overflow-visible text-foreground"
        >
          <path
            d="M10 6.5v19.25C10 31.8 13.6 35 20.3 35H31"
            fill="none"
            stroke="currentColor"
            strokeWidth="3.8"
            strokeLinecap="round"
            strokeLinejoin="round"
          />
          <text
            x="19"
            y="27.5"
            fill="currentColor"
            style={{
              fontFamily: "var(--font-loreline-story)",
              fontSize: 28,
              fontStyle: "italic",
              fontWeight: 600,
              letterSpacing: "-1.35px",
            }}
          >
            oreline
          </text>
          <path
            d="M27 35H122"
            fill="none"
            stroke="currentColor"
            strokeWidth="3.8"
            strokeLinecap="round"
          />
          <circle cx="130" cy="35" r="3.4" fill="var(--sky)" />
        </svg>
      )}
    </Link>
  );
};

export function LorelineMark({ className }: { className?: string }) {
  return (
    <span
      aria-hidden="true"
      className={cn(
        "relative inline-grid size-9 place-items-center rounded-[0.7rem] bg-foreground text-background",
        className,
      )}
    >
      <span className="-translate-y-px font-semibold font-story text-[1.65rem] italic leading-none">
        L
      </span>
      <span className="absolute bottom-[0.42rem] left-[0.48rem] h-[1.5px] w-[1.05rem] bg-brand" />
      <span className="absolute right-[0.42rem] bottom-[0.31rem] size-1 rounded-full bg-sky" />
    </span>
  );
}
