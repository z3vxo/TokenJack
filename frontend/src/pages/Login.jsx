import { useState } from "react";
import { useNavigate } from "react-router-dom";
import {
  CircleNotch,
  Eye,
  EyeSlash,
  WarningCircle,
} from "@phosphor-icons/react";
import BrandMark from "../components/BrandMark.jsx";
import ThemeToggle from "../components/ThemeToggle.jsx";
import { login } from "../lib/api.js";
import { setTokens } from "../lib/auth.js";

const FIELD =
  "h-11 w-full rounded-control border bg-card px-3.5 text-sm text-foreground " +
  "placeholder:text-placeholder outline-none transition-colors " +
  "focus-visible:border-primary focus-visible:ring-[3px] " +
  "focus-visible:ring-[color:var(--focus-ring-color)] " +
  "disabled:cursor-not-allowed disabled:bg-muted disabled:opacity-70";

export default function Login() {
  const navigate = useNavigate();
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const [error, setError] = useState("");
  const [submitting, setSubmitting] = useState(false);

  async function handleSubmit(event) {
    event.preventDefault();
    if (submitting) return;

    setError("");
    setSubmitting(true);
    try {
      const tokens = await login(username.trim(), password);
      setTokens(tokens);
      navigate("/", { replace: true });
    } catch (err) {
      setError(err?.message || "Unable to sign in.");
      setSubmitting(false);
    }
  }

  const hasError = Boolean(error);

  return (
    <div className="flex min-h-[100dvh] flex-col bg-background">
      <header className="relative flex items-center justify-center px-6 pt-10 sm:pt-14">
        <div className="flex items-center gap-2.5">
          <BrandMark />
          <span className="text-[0.95rem] font-semibold tracking-[-0.01em]">
            TokenJack
          </span>
        </div>
        <ThemeToggle className="absolute right-6 top-10 sm:top-14" />
      </header>

      <main className="flex flex-1 items-center justify-center px-6 py-10">
        <div className="w-full max-w-[400px]">
          <div className="mb-8">
            <h1 className="text-[1.75rem] font-semibold leading-tight tracking-[-0.02em]">
              Sign in
            </h1>
            <p className="mt-2 text-sm leading-6 text-muted-foreground">
              Enter your operator credentials to access the admin console.
            </p>
          </div>

          <form
            onSubmit={handleSubmit}
            noValidate
            aria-busy={submitting}
            className="space-y-5"
          >
            {hasError && (
              <div
                id="login-error"
                role="alert"
                className="flex items-start gap-2.5 rounded-control border border-destructive-border bg-destructive-surface px-3.5 py-3 text-sm text-destructive"
              >
                <WarningCircle
                  size={18}
                  weight="fill"
                  aria-hidden="true"
                  className="mt-0.5 shrink-0"
                />
                <span className="leading-5">{error}</span>
              </div>
            )}

            <div className="space-y-2">
              <label
                htmlFor="username"
                className="block text-sm font-medium text-foreground"
              >
                Username
              </label>
              <input
                id="username"
                name="username"
                type="text"
                autoComplete="username"
                autoCapitalize="none"
                autoCorrect="off"
                spellCheck={false}
                required
                autoFocus
                disabled={submitting}
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                aria-invalid={hasError || undefined}
                aria-describedby={hasError ? "login-error" : undefined}
                placeholder="admin"
                className={`${FIELD} ${
                  hasError ? "border-destructive" : "border-input"
                }`}
              />
            </div>

            <div className="space-y-2">
              <label
                htmlFor="password"
                className="block text-sm font-medium text-foreground"
              >
                Password
              </label>
              <div className="relative">
                <input
                  id="password"
                  name="password"
                  type={showPassword ? "text" : "password"}
                  autoComplete="current-password"
                  required
                  disabled={submitting}
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  aria-invalid={hasError || undefined}
                  aria-describedby={hasError ? "login-error" : undefined}
                  placeholder="••••••••"
                  className={`${FIELD} pr-11 ${
                    hasError ? "border-destructive" : "border-input"
                  }`}
                />
                <button
                  type="button"
                  onClick={() => setShowPassword((v) => !v)}
                  aria-label={showPassword ? "Hide password" : "Show password"}
                  aria-pressed={showPassword}
                  tabIndex={submitting ? -1 : 0}
                  className="absolute inset-y-0 right-0 flex w-11 cursor-pointer items-center justify-center rounded-control text-muted-foreground transition-colors hover:text-foreground focus-visible:text-foreground"
                >
                  {showPassword ? (
                    <EyeSlash size={18} aria-hidden="true" />
                  ) : (
                    <Eye size={18} aria-hidden="true" />
                  )}
                </button>
              </div>
            </div>

            <button
              type="submit"
              disabled={submitting}
              className="flex h-11 w-full cursor-pointer items-center justify-center gap-2 rounded-control bg-primary text-sm font-medium text-on-primary transition-colors hover:bg-primary-hover focus-visible:ring-[3px] focus-visible:ring-[color:var(--focus-ring-color)] disabled:cursor-not-allowed disabled:opacity-60"
            >
              {submitting && (
                <CircleNotch
                  size={17}
                  className="animate-spin"
                  aria-hidden="true"
                />
              )}
              {submitting ? "Signing in…" : "Sign in"}
            </button>
          </form>
        </div>
      </main>

      <footer className="px-6 pb-8 text-center">
        <p className="text-xs text-muted-foreground">
          TokenJack · Authorized operators only
        </p>
      </footer>
    </div>
  );
}
