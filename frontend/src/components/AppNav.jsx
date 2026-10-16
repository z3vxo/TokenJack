import { NavLink } from "react-router-dom";

const TABS = [
  { to: "/overview", label: "Overview" },
  { to: "/campaigns", label: "Campaigns" },
  { to: "/automation", label: "Automation" },
  { to: "/tokens", label: "Tokens" },
  { to: "/settings", label: "Settings" },
];

export default function AppNav() {
  return (
    <nav
      aria-label="Sections"
      className="overflow-x-auto border-b border-border bg-card"
    >
      <ul className="mx-auto flex w-max items-center gap-10 px-6">
        {TABS.map((tab) => (
          <li key={tab.to} className="shrink-0">
            <NavLink
              to={tab.to}
              className={({ isActive }) =>
                `-mb-px inline-flex h-11 cursor-pointer items-center border-b-2 text-sm transition-colors ${
                  isActive
                    ? "border-primary font-medium text-foreground"
                    : "border-transparent text-muted-foreground hover:text-foreground"
                }`
              }
            >
              {tab.label}
            </NavLink>
          </li>
        ))}
      </ul>
    </nav>
  );
}
