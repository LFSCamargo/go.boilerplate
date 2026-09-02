import type { TailwindConfig } from "react-email";
import plugin from "tailwindcss/plugin";

/** Go Boilerplate emails — Protocol tokens from Figma SaaS Email Templates. */
const colors = {
  bg: "#131313",
  "bg-2": "#212121",
  fg: "#FFFFFF",
  "fg-2": "#C4C4C4",
  "fg-3": "#818181",
  "fg-inverted": "#131313",
  stroke: "#2B2B2B",
  brand: "#FFFFFF",
  surface: "#FFFFFF",
  ink: "#131313",
} as const;

const fontScale = {
  11: {
    fontSize: "11px",
    lineHeight: "1.5",
    letterSpacing: "-0.033px",
    fontWeight: "500",
  },
  13: {
    fontSize: "13px",
    lineHeight: "1.5",
    letterSpacing: "-0.039px",
    fontWeight: "500",
  },
  15: {
    fontSize: "15px",
    lineHeight: "1.5",
    letterSpacing: "-0.075px",
    fontWeight: "400",
  },
  "15-medium": {
    fontSize: "15px",
    lineHeight: "1.5",
    letterSpacing: "-0.075px",
    fontWeight: "500",
  },
  17: {
    fontSize: "17px",
    lineHeight: "1.5",
    letterSpacing: "-0.51px",
    fontWeight: "500",
  },
  40: {
    fontSize: "40px",
    lineHeight: "1.05",
    letterSpacing: "-0.8px",
    fontWeight: "500",
  },
  56: {
    fontSize: "56px",
    lineHeight: "1.05",
    letterSpacing: "-1.12px",
    fontWeight: "500",
  },
} as const;

export const goBoilerplateTailwindConfig: TailwindConfig = {
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
      colors,
      fontFamily: {
        sans: ["Inter", "Arial", "sans-serif"],
        condensed: [
          "IBM Plex Sans Condensed",
          "Arial Narrow",
          "sans-serif",
        ],
      },
    },
  },
};
