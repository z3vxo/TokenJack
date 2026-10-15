import { Outlet, useNavigate } from "react-router-dom";
import { SignOut } from "@phosphor-icons/react";
import BrandMark from "./BrandMark.jsx";
import CampaignDropdown from "./CampaignDropdown.jsx";
import AppNav from "./AppNav.jsx";
import ThemeToggle from "./ThemeToggle.jsx";
import { clearTokens } from "../lib/auth.js";

export default function AppLayout() {
  const navigate = useNavigate();

  function signOut() {
    clearTokens();
    navigate("/login", { replace: true });
  }

  return (
    <div className="flex min-h-[100dvh] flex-col bg-background">
      <header className="flex items-center justify-between border-b border-border bg-card px-6 py-3.5">
        <div className="flex items-center gap-2.5">
          <BrandMark size={28} />
          <span className="text-sm font-semibold tracking-[-0.01em]">
            TokenJack
          </span>
        </div>
        <div className="flex items-center gap-3">
          <CampaignDropdown />
          <ThemeToggle />
          <button
            type="button"
            onClick={signOut}
            className="flex cursor-pointer items-center gap-2 rounded-control border border-border bg-card px-3 py-2 text-sm font-medium text-foreground transition-colors hover:bg-muted focus-visible:ring-[3px] focus-visible:ring-[color:var(--focus-ring-color)]"
          >
            <SignOut size={16} aria-hidden="true" />
            Sign out
          </button>
        </div>
      </header>

      <AppNav />

      <main className="flex flex-1 px-6 py-8">
        <Outlet />
      </main>
    </div>
  );
}
