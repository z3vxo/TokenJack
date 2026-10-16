import { useCallback, useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import {
  CircleNotch,
  Cube,
  CursorClick,
  Key,
  Stack,
  WarningCircle,
} from "@phosphor-icons/react";
import { useCampaign } from "../context/CampaignContext.jsx";
import { useCampaignEvent } from "../context/EventsContext.jsx";
import { getCampaignMeta } from "../lib/api.js";
import { clearTokens } from "../lib/auth.js";

export default function Overview() {
  const navigate = useNavigate();
  const { selected } = useCampaign();
  const [meta, setMeta] = useState(null);
  const [status, setStatus] = useState("idle");

  const load = useCallback(
    async ({ silent = false } = {}) => {
      const id = selected?.uniqueURI;
      if (!id) {
        setMeta(null);
        setStatus("idle");
        return;
      }

      if (!silent) setStatus("loading");
      try {
        const data = await getCampaignMeta(id);
        setMeta(data);
        setStatus("ready");
      } catch (err) {
        if (err?.status === 401) {
          clearTokens();
          navigate("/login", { replace: true });
          return;
        }
        if (!silent) setStatus("error");
      }
    },
    [selected?.uniqueURI, navigate],
  );

  useEffect(() => {
    load();
  }, [load]);

  useCampaignEvent((event) => {
    if (event?.type === "token") load({ silent: true });
  });

  if (!selected) return <EmptyState />;
  if (status === "loading") return <LoadingState />;
  if (status === "error") return <ErrorState />;
  if (!meta) return <LoadingState />;

  return <MetaView meta={meta} />;
}

function EmptyState() {
  return (
    <div className="m-auto max-w-sm text-center">
      <h1 className="text-xl font-semibold tracking-[-0.02em]">
        No campaign selected
      </h1>
      <p className="mt-2 text-sm leading-6 text-muted-foreground">
        Please choose a campaign in the top right.
      </p>
    </div>
  );
}

function LoadingState() {
  return (
    <div
      role="status"
      aria-live="polite"
      className="m-auto flex items-center gap-2 text-muted-foreground"
    >
      <CircleNotch size={18} className="animate-spin" aria-hidden="true" />
      <span className="text-sm">Loading overview…</span>
    </div>
  );
}

function ErrorState() {
  return (
    <div
      role="alert"
      className="m-auto flex max-w-sm items-center gap-2.5 rounded-control border border-destructive-border bg-destructive-surface px-3.5 py-3 text-sm text-destructive"
    >
      <WarningCircle
        size={18}
        weight="fill"
        aria-hidden="true"
        className="shrink-0"
      />
      <span>Failed to load the campaign overview.</span>
    </div>
  );
}

function MetaView({ meta }) {
  const clicks = meta.clicks ?? 0;
  const tokens = meta.tokens_total ?? 0;
  const scopes = Array.isArray(meta.scopes) ? meta.scopes : [];

  return (
    <div className="mx-auto w-full max-w-5xl">
      <header className="mb-6">
        <h1 className="text-xl font-semibold tracking-[-0.02em]">
          {meta.name}
        </h1>
        <p className="mt-1 font-mono text-xs text-muted-foreground">
          {meta.unique_str}
        </p>
      </header>

      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
        <Stat icon={CursorClick} label="Clicks" value={clicks} />
        <Stat icon={Key} label="Tokens issued" value={tokens} />
        <Stat icon={Stack} label="Scopes" value={scopes.length} />
        <Stat icon={Cube} label="Resource" value={meta.resource} mono />
      </div>

      <section className="mt-6 rounded-control border border-border bg-card p-5">
        <h2 className="text-sm font-medium text-foreground">Engagement</h2>
        <p className="mt-1 text-xs text-muted-foreground">
          Clicks and tokens issued for this campaign.
        </p>
        <div className="mt-5">
          <MetricsChart clicks={clicks} tokens={tokens} />
        </div>
      </section>

      <section className="mt-6 rounded-control border border-border bg-card p-5">
        <h2 className="text-sm font-medium text-foreground">Scopes</h2>
        {scopes.length === 0 ? (
          <p className="mt-2 text-sm text-muted-foreground">
            No scopes configured.
          </p>
        ) : (
          <ul className="mt-3 flex flex-wrap gap-2">
            {scopes.map((scope) => (
              <li
                key={scope}
                className="rounded-control border border-border px-2.5 py-1 font-mono text-xs text-muted-foreground"
              >
                {scope}
              </li>
            ))}
          </ul>
        )}
      </section>
    </div>
  );
}

function Stat({ icon: Icon, label, value, mono = false }) {
  return (
    <div className="rounded-control border border-border bg-card p-4">
      <div className="flex items-center gap-2 text-muted-foreground">
        <Icon size={15} aria-hidden="true" />
        <span className="text-xs font-medium">{label}</span>
      </div>
      <p
        className={`mt-2 font-semibold tracking-[-0.02em] text-foreground ${
          mono ? "truncate font-mono text-lg" : "text-2xl tabular-nums"
        }`}
      >
        {value}
      </p>
    </div>
  );
}

function MetricsChart({ clicks, tokens }) {
  const max = Math.max(clicks, tokens, 1);
  const rows = [
    { label: "Clicks", value: clicks, bar: "bg-primary" },
    { label: "Tokens issued", value: tokens, bar: "bg-success" },
  ];

  return (
    <div className="space-y-5">
      {rows.map((row) => (
        <div key={row.label}>
          <div className="mb-2 flex items-baseline justify-between text-sm">
            <span className="text-muted-foreground">{row.label}</span>
            <span className="font-medium tabular-nums text-foreground">
              {row.value}
            </span>
          </div>
          <div className="h-2 w-full bg-muted" aria-hidden="true">
            <div
              className={`h-full ${row.bar} transition-[width] duration-300`}
              style={{ width: `${(row.value / max) * 100}%` }}
            />
          </div>
        </div>
      ))}
    </div>
  );
}
