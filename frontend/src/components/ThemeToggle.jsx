import { useState } from "react";
import { Moon, Sun } from "@phosphor-icons/react";
import { getTheme, setTheme } from "../lib/theme.js";

export default function ThemeToggle({ className = "" }) {
  const [theme, setThemeState] = useState(
    () => document.documentElement.dataset.theme || getTheme(),
  );

  function toggle() {
    const next = theme === "dark" ? "light" : "dark";
    setThemeState(next);
    setTheme(next);
  }

  const isDark = theme === "dark";

  return (
    <button
      type="button"
      onClick={toggle}
      aria-label={isDark ? "Switch to light mode" : "Switch to dark mode"}
      title={isDark ? "Light mode" : "Dark mode"}
      className={`flex h-9 w-9 cursor-pointer items-center justify-center rounded-control border border-border bg-card text-muted-foreground transition-colors hover:bg-muted hover:text-foreground focus-visible:ring-[3px] focus-visible:ring-[color:var(--focus-ring-color)] ${className}`}
    >
      {isDark ? (
        <Sun size={16} aria-hidden="true" />
      ) : (
        <Moon size={16} aria-hidden="true" />
      )}
    </button>
  );
}
