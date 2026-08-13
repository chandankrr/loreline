import type { Metadata } from "next";
import { IBM_Plex_Mono, Manrope, Newsreader } from "next/font/google";

import { cn } from "@loreline/ui/lib/utils";
import "@loreline/ui/globals.css";

import { TooltipProvider } from "@loreline/ui/components/tooltip";

const sans = Manrope({ subsets: ["latin"], variable: "--font-loreline-sans" });
const story = Newsreader({
  subsets: ["latin"],
  weight: ["500", "600"],
  style: ["normal", "italic"],
  variable: "--font-loreline-story",
});
const mono = IBM_Plex_Mono({
  subsets: ["latin"],
  weight: ["400", "500"],
  variable: "--font-mono",
});

export const metadata: Metadata = {
  title: "Loreline",
  description:
    "A realtime AI reading companion that sees the page with you, answers by voice, and turns difficult ideas into vivid understanding.",
};

export default function RootLayout({ children }: LayoutProps<"/">) {
  return (
    <html
      lang="en"
      suppressHydrationWarning
      data-scroll-behavior="smooth"
      className={cn(
        "h-full antialiased",
        sans.variable,
        story.variable,
        mono.variable,
      )}
    >
      <body className="min-h-full">
        <TooltipProvider>{children}</TooltipProvider>
      </body>
    </html>
  );
}
