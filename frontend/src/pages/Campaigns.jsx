import { useRef, useState } from "react";
import { useNavigate } from "react-router-dom";
import { CircleNotch, Plus, X } from "@phosphor-icons/react";
import Modal from "../components/Modal.jsx";
import { useToast } from "../components/Toast.jsx";
import { useCampaign } from "../context/CampaignContext.jsx";
import { createCampaign } from "../lib/api.js";
import { clearTokens } from "../lib/auth.js";

const FIELD =
  "h-11 w-full rounded-control border border-input bg-card px-3.5 text-sm text-foreground " +
  "placeholder:text-placeholder outline-none transition-colors " +
  "focus-visible:border-primary focus-visible:ring-[3px] " +
  "focus-visible:ring-[color:var(--focus-ring-color)]";

export default function Campaigns() {
  const navigate = useNavigate();
  const { campaigns, status, reload } = useCampaign();
  const toast = useToast();
  const [open, setOpen] = useState(false);

  function handleUnauthorized() {
    clearTokens();
    navigate("/login", { replace: true });
  }

  return (
    <div className="mx-auto w-full max-w-5xl">
      <header className="mb-6 flex items-center justify-between gap-4">
        <div>
          <h1 className="text-xl font-semibold tracking-[-0.02em]">
            Campaigns
          </h1>
          <p className="mt-1 text-sm text-muted-foreground">
            {campaigns.length}{" "}
            {campaigns.length === 1 ? "campaign" : "campaigns"}
          </p>
        </div>
        <button
          type="button"
          onClick={() => setOpen(true)}
          className="flex h-11 cursor-pointer items-center gap-2 rounded-control bg-primary px-4 text-sm font-medium text-on-primary transition-colors hover:bg-primary-hover focus-visible:ring-[3px] focus-visible:ring-[color:var(--focus-ring-color)]"
        >
          <Plus size={16} weight="bold" aria-hidden="true" />
          New
        </button>
      </header>

      <CampaignList campaigns={campaigns} status={status} />

      <Modal
        open={open}
        onClose={() => setOpen(false)}
        labelledBy="new-campaign-title"
      >
        <NewCampaignForm
          onCancel={() => setOpen(false)}
          onCreated={() => {
            setOpen(false);
            toast.success("Campaign created.");
            reload();
          }}
          onError={(message) => toast.error(message)}
          onUnauthorized={handleUnauthorized}
        />
      </Modal>
    </div>
  );
}

function CampaignList({ campaigns, status }) {
  if (status === "loading" && campaigns.length === 0) {
    return (
      <div
        role="status"
        aria-live="polite"
        className="flex items-center gap-2 text-sm text-muted-foreground"
      >
        <CircleNotch size={16} className="animate-spin" aria-hidden="true" />
        Loading campaigns…
      </div>
    );
  }

  if (campaigns.length === 0) {
    return (
      <div className="rounded-control border border-border bg-card px-4 py-10 text-center">
        <p className="text-sm text-muted-foreground">
          No campaigns yet. Create one to get started.
        </p>
      </div>
    );
  }

  return (
    <div className="overflow-hidden rounded-control border border-border bg-card">
      <div className="hidden grid-cols-[1.5fr_1.5fr_1fr] gap-4 border-b border-border px-4 py-2.5 text-xs font-medium text-muted-foreground sm:grid">
        <span>Name</span>
        <span>Tenant ID</span>
        <span className="text-right">Unique string</span>
      </div>
      <ul className="divide-y divide-border">
        {campaigns.map((campaign) => (
          <li
            key={campaign.uniqueURI}
            className="grid gap-1 px-4 py-3 sm:grid-cols-[1.5fr_1.5fr_1fr] sm:items-center sm:gap-4"
          >
            <span className="truncate text-sm font-medium text-foreground">
              {campaign.name}
            </span>
            <span className="truncate font-mono text-xs text-muted-foreground">
              {campaign.tenant_id}
            </span>
            <span className="truncate font-mono text-xs text-muted-foreground sm:text-right">
              {campaign.uniqueURI}
            </span>
          </li>
        ))}
      </ul>
    </div>
  );
}

function NewCampaignForm({ onCancel, onCreated, onError, onUnauthorized }) {
  const [name, setName] = useState("");
  const [tenantId, setTenantId] = useState("");
  const [resource, setResource] = useState("graph");
  const [useDefaultScope, setUseDefaultScope] = useState(true);
  const [customScopes, setCustomScopes] = useState([{ id: 0, value: "" }]);
  const scopeId = useRef(0);
  const [submitting, setSubmitting] = useState(false);

  function toggleDefault(checked) {
    setUseDefaultScope(checked);
    if (!checked && customScopes.length === 0) {
      setCustomScopes([{ id: (scopeId.current += 1), value: "" }]);
    }
  }

  function updateScope(id, value) {
    setCustomScopes((list) =>
      list.map((scope) => (scope.id === id ? { ...scope, value } : scope)),
    );
  }

  function addScope() {
    setCustomScopes((list) => [...list, { id: (scopeId.current += 1), value: "" }]);
  }

  function removeScope(id) {
    setCustomScopes((list) => list.filter((scope) => scope.id !== id));
  }

  async function handleSubmit(event) {
    event.preventDefault();
    if (submitting) return;

    const scopes = useDefaultScope
      ? [".default"]
      : customScopes.map((scope) => scope.value.trim()).filter(Boolean);

    if (scopes.length === 0) {
      onError("Enter at least one scope.");
      return;
    }

    setSubmitting(true);
    try {
      await createCampaign({
        name: name.trim(),
        tenantId: tenantId.trim(),
        resource,
        scopes,
      });
      onCreated();
    } catch (err) {
      if (err?.status === 401) {
        onUnauthorized();
        return;
      }
      onError(err?.message || "Could not create the campaign.");
      setSubmitting(false);
    }
  }

  return (
    <div className="w-full p-5">
      <h2
        id="new-campaign-title"
        className="text-base font-semibold tracking-[-0.01em]"
      >
        New campaign
      </h2>
      <p className="mt-1 text-sm text-muted-foreground">
        Register a campaign to issue tokens for a tenant.
      </p>

      <form onSubmit={handleSubmit} className="mt-5 space-y-4">
        <div className="space-y-2">
          <label
            htmlFor="campaign-name"
            className="block text-sm font-medium text-foreground"
          >
            Name
          </label>
          <input
            id="campaign-name"
            type="text"
            required
            autoFocus
            disabled={submitting}
            value={name}
            onChange={(e) => setName(e.target.value)}
            placeholder="Acme rollout"
            className={FIELD}
          />
        </div>

        <div className="space-y-2">
          <label
            htmlFor="campaign-tenant"
            className="block text-sm font-medium text-foreground"
          >
            Tenant ID
          </label>
          <input
            id="campaign-tenant"
            type="text"
            required
            spellCheck={false}
            autoCapitalize="none"
            autoCorrect="off"
            disabled={submitting}
            value={tenantId}
            onChange={(e) => setTenantId(e.target.value)}
            placeholder="00000000-0000-0000-0000-000000000000"
            className={`${FIELD} font-mono`}
          />
        </div>

        <div className="space-y-2">
          <label
            htmlFor="campaign-resource"
            className="block text-sm font-medium text-foreground"
          >
            Resource
          </label>
          <select
            id="campaign-resource"
            disabled={submitting}
            value={resource}
            onChange={(e) => setResource(e.target.value)}
            className={`${FIELD} cursor-pointer appearance-none font-mono`}
          >
            <option value="graph">graph</option>
          </select>
        </div>

        <fieldset className="space-y-3">
          <legend className="text-sm font-medium text-foreground">Scopes</legend>
          <label
            htmlFor="scope-default"
            className="flex cursor-pointer items-center gap-2.5 rounded-control border border-border px-3 py-2.5 text-sm transition-colors hover:bg-muted"
          >
            <input
              id="scope-default"
              type="checkbox"
              checked={useDefaultScope}
              disabled={submitting}
              onChange={(e) => toggleDefault(e.target.checked)}
              className="h-4 w-4 cursor-pointer accent-primary"
            />
            <span className="font-mono text-foreground">.default</span>
          </label>

          {!useDefaultScope && (
            <div className="space-y-2">
              {customScopes.map((scope, index) => (
                <div key={scope.id} className="flex items-center gap-2">
                  <input
                    type="text"
                    value={scope.value}
                    disabled={submitting}
                    spellCheck={false}
                    autoCapitalize="none"
                    autoCorrect="off"
                    onChange={(e) => updateScope(scope.id, e.target.value)}
                    placeholder="e.g. User.Read"
                    aria-label={`Scope ${index + 1}`}
                    className={`${FIELD} font-mono`}
                  />
                  {customScopes.length > 1 && (
                    <button
                      type="button"
                      onClick={() => removeScope(scope.id)}
                      disabled={submitting}
                      aria-label={`Remove scope ${index + 1}`}
                      className="flex h-11 w-11 shrink-0 cursor-pointer items-center justify-center rounded-control border border-border text-muted-foreground transition-colors hover:bg-muted hover:text-foreground focus-visible:ring-[3px] focus-visible:ring-[color:var(--focus-ring-color)] disabled:cursor-not-allowed disabled:opacity-60"
                    >
                      <X size={15} aria-hidden="true" />
                    </button>
                  )}
                </div>
              ))}
              <button
                type="button"
                onClick={addScope}
                disabled={submitting}
                className="inline-flex cursor-pointer items-center gap-1.5 rounded-control text-sm font-medium text-primary transition-colors hover:text-primary-hover focus-visible:ring-[3px] focus-visible:ring-[color:var(--focus-ring-color)] disabled:cursor-not-allowed disabled:opacity-60"
              >
                <Plus size={15} weight="bold" aria-hidden="true" />
                Add scope
              </button>
            </div>
          )}
        </fieldset>

        <div className="flex justify-end gap-2 pt-1">
          <button
            type="button"
            onClick={onCancel}
            disabled={submitting}
            className="flex h-11 cursor-pointer items-center rounded-control border border-border bg-card px-4 text-sm font-medium text-foreground transition-colors hover:bg-muted focus-visible:ring-[3px] focus-visible:ring-[color:var(--focus-ring-color)] disabled:cursor-not-allowed disabled:opacity-60"
          >
            Cancel
          </button>
          <button
            type="submit"
            disabled={submitting}
            className="flex h-11 cursor-pointer items-center justify-center gap-2 rounded-control bg-primary px-4 text-sm font-medium text-on-primary transition-colors hover:bg-primary-hover focus-visible:ring-[3px] focus-visible:ring-[color:var(--focus-ring-color)] disabled:cursor-not-allowed disabled:opacity-60"
          >
            {submitting && (
              <CircleNotch
                size={16}
                className="animate-spin"
                aria-hidden="true"
              />
            )}
            {submitting ? "Creating…" : "Create"}
          </button>
        </div>
      </form>
    </div>
  );
}
