import type { Metadata } from "next";
import { IBM_Plex_Mono, Manrope, Newsreader } from "next/font/google";

import { cn } from "@loreline/ui/lib/utils";
import "@loreline/ui/globals.css";

import { Toaster } from "@loreline/ui/components/toast";
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
  metadataBase: new URL(
    process.env.NEXT_PUBLIC_APP_URL ?? "http://localhost:3000",
  ),
  title: "Loreline",
  description:
    "A realtime AI reading companion that sees the page with you, answers by voice, and turns difficult ideas into vivid understanding.",
  applicationName: "Loreline",
  icons: {
    icon: [{ url: "/icon.svg", type: "image/svg+xml" }],
    shortcut: "/icon.svg",
  },
  openGraph: {
    title: "Loreline — Read past the words",
    description: "Your books, made conversational and visual.",
    type: "website",
  },
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
        <Toaster />
      </body>
    </html>
  );
}
