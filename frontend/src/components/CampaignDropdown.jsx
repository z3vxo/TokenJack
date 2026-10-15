import { useEffect, useRef, useState } from "react";
import { Buildings, CaretDown, Check } from "@phosphor-icons/react";
import { useCampaign } from "../context/CampaignContext.jsx";

export default function CampaignDropdown() {
  const { campaigns, selected, select, status } = useCampaign();
  const [open, setOpen] = useState(false);
  const ref = useRef(null);

  useEffect(() => {
    if (!open) return;

    function handlePointer(event) {
      if (ref.current && !ref.current.contains(event.target)) setOpen(false);
    }
    function handleKey(event) {
      if (event.key === "Escape") setOpen(false);
    }

    document.addEventListener("mousedown", handlePointer);
    document.addEventListener("keydown", handleKey);
    return () => {
      document.removeEventListener("mousedown", handlePointer);
      document.removeEventListener("keydown", handleKey);
    };
  }, [open]);

  const loading = status === "loading";
  const label = loading ? "Loading…" : selected?.name || "Select campaign";

  return (
    <div className="relative" ref={ref}>
      <button
        type="button"
        onClick={() => setOpen((v) => !v)}
        disabled={loading}
        aria-haspopup="listbox"
        aria-expanded={open}
        aria-label="Current campaign"
        className="flex w-64 cursor-pointer items-center gap-2 rounded-control border border-border bg-card px-3 py-2 text-sm transition-colors hover:bg-muted focus-visible:ring-[3px] focus-visible:ring-[color:var(--focus-ring-color)] disabled:cursor-default disabled:opacity-70"
      >
        <Buildings
          size={16}
          className="shrink-0 text-muted-foreground"
          aria-hidden="true"
        />
        <span
          className={`flex-1 truncate text-left ${
            selected ? "font-medium text-foreground" : "text-muted-foreground"
          }`}
        >
          {label}
        </span>
        <CaretDown
          size={14}
          className={`shrink-0 text-muted-foreground transition-transform duration-200 ${
            open ? "rotate-180" : ""
          }`}
          aria-hidden="true"
        />
      </button>

      {open && (
        <div
          role="listbox"
          aria-label="Campaigns"
          className="absolute right-0 z-30 mt-2 w-64 overflow-hidden rounded-control border border-border bg-card shadow-menu"
        >
          {campaigns.length === 0 ? (
            <p className="px-3.5 py-3 text-sm text-muted-foreground">
              No campaigns found.
            </p>
          ) : (
            <ul className="max-h-72 overflow-y-auto py-1">
              {campaigns.map((campaign) => {
                const active = selected?.uniqueURI === campaign.uniqueURI;
                return (
                  <li key={campaign.uniqueURI}>
                    <button
                      type="button"
                      role="option"
                      aria-selected={active}
                      onClick={() => {
                        select(campaign);
                        setOpen(false);
                      }}
                      className="flex w-full cursor-pointer items-center gap-2 px-3.5 py-2 text-left text-sm transition-colors hover:bg-muted focus-visible:bg-muted"
                    >
                      <span
                        className={`flex-1 truncate ${
                          active ? "font-medium text-foreground" : "text-foreground"
                        }`}
                      >
                        {campaign.name}
                      </span>
                      {active && (
                        <Check
                          size={15}
                          weight="bold"
                          className="shrink-0 text-primary"
                          aria-hidden="true"
                        />
                      )}
                    </button>
                  </li>
                );
              })}
            </ul>
          )}
        </div>
      )}
    </div>
  );
}
