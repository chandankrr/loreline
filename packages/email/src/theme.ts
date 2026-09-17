import type { TailwindConfig } from "react-email";
import plugin from "tailwindcss/plugin";

export const loreline = {
  paper: "#FAFAF1",
  "paper-deep": "#E3E3D0",
  surface: "#FFFFFF",
  "surface-muted": "#F3F2E2",
  ink: "#13150C", // primary text
  "ink-soft": "#4B4E3B", // secondary text
  "ink-faint": "#7C806B", // tertiary / caption text
  border: "#DEDCC6",
  brand: "#DFE044",
  "brand-ink": "#2F4300",
  "brand-soft": "#EFF0C4",
  accent: "#FFD7EE",
  "accent-Ink": "#670242",
  sage: "#59985B",
  "sage-Soft": "#E7F3E7",
  gold: "#D9B255",
} as const;

const fontScale = {
  12: { fontSize: "12px", lineHeight: "1.6", letterSpacing: "0px" },
  13: { fontSize: "13px", lineHeight: "1.55", letterSpacing: "0px" },
  14: { fontSize: "14px", lineHeight: "1.6", letterSpacing: "0px" },
  16: { fontSize: "16px", lineHeight: "1.6", letterSpacing: "-0.1px" },
  18: { fontSize: "18px", lineHeight: "1.4", letterSpacing: "-0.15px" },
  22: { fontSize: "22px", lineHeight: "1.3", letterSpacing: "-0.2px" },
  30: { fontSize: "30px", lineHeight: "1.2", letterSpacing: "-0.4px" },
  40: { fontSize: "40px", lineHeight: "1.08", letterSpacing: "-0.6px" },
} as const;

export const lorelineTailwindConfig: TailwindConfig = {
  plugins: [
    plugin(({ addUtilities, addVariant }) => {
      addVariant("mobile", "@media (max-width: 600px)");
      const utilities: Record<string, Record<string, string>> = {};
      for (const [step, token] of Object.entries(fontScale)) {
        utilities[`.font-${step}`] = token;
      }
      addUtilities(utilities);
    }),
  ],
  theme: {
    extend: {
      colors: loreline,
      borderRadius: {
        sm: "6px",
        DEFAULT: "10px",
        lg: "14px",
        xl: "20px",
        full: "999px",
      },
      boxShadow: {
        card: "0px 18px 44px 0px rgba(19, 21, 12, 0.06), 0px 4px 10px 0px rgba(19, 21, 12, 0.04)",
      },
      fontFamily: {
        // Body / UI text
        sans: ["Manrope", "Helvetica Neue", "Arial", "sans-serif"],
        // Headlines and epigraph-style quotes — set in italic in components
        display: ["Newsreader", "Georgia", "Times New Roman", "serif"],
      },
    },
  },
};

export default lorelineTailwindConfig;
