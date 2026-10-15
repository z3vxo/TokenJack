import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useState,
} from "react";
import { useNavigate } from "react-router-dom";
import { getCampaigns } from "../lib/api.js";
import { clearTokens } from "../lib/auth.js";
import { loadCampaign, saveCampaign } from "../lib/campaigns.js";

const CampaignContext = createContext(null);

export function CampaignProvider({ children }) {
  const navigate = useNavigate();
  const [campaigns, setCampaigns] = useState([]);
  const [selected, setSelected] = useState(loadCampaign);
  const [status, setStatus] = useState("loading");

  useEffect(() => {
    let cancelled = false;

    (async () => {
      try {
        const data = await getCampaigns();
        if (cancelled) return;

        const list = Array.isArray(data?.campaigns) ? data.campaigns : [];
        setCampaigns(list);
        setStatus("ready");

        setSelected((current) => {
          if (current && !list.some((c) => c.uniqueURI === current.uniqueURI)) {
            saveCampaign(null);
            return null;
          }
          return current;
        });
      } catch (err) {
        if (cancelled) return;

        if (err?.status === 401) {
          clearTokens();
          navigate("/login", { replace: true });
          return;
        }

        setStatus("error");
      }
    })();

    return () => {
      cancelled = true;
    };
  }, [navigate]);

  const select = useCallback((campaign) => {
    setSelected(campaign);
    saveCampaign(campaign);
  }, []);

  return (
    <CampaignContext.Provider value={{ campaigns, selected, select, status }}>
      {children}
    </CampaignContext.Provider>
  );
}

export function useCampaign() {
  const ctx = useContext(CampaignContext);
  if (!ctx) {
    throw new Error("useCampaign must be used within a CampaignProvider");
  }
  return ctx;
}
