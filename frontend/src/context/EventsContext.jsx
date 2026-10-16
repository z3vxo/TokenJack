import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useRef,
  useState,
} from "react";
import { getAccess } from "../lib/auth.js";
import { campaignEventsUrl } from "../lib/events.js";
import { useCampaign } from "./CampaignContext.jsx";

const EventsContext = createContext(null);

export function EventsProvider({ children }) {
  const { selected } = useCampaign();
  const [status, setStatus] = useState("idle");
  const listeners = useRef(new Set());

  useEffect(() => {
    const id = selected?.uniqueURI;
    if (!id) {
      setStatus("idle");
      return;
    }

    const token = getAccess();
    if (!token) {
      setStatus("closed");
      return;
    }

    setStatus("connecting");
    const source = new EventSource(campaignEventsUrl(id, token));
    source.onopen = () => setStatus("open");
    source.onerror = () => setStatus("closed");
    source.onmessage = (event) => {
      let data = event.data;
      try {
        data = JSON.parse(event.data);
      } catch {
        // keep raw string payloads as-is
      }
      for (const listener of listeners.current) listener(data, event);
    };

    return () => {
      source.close();
    };
  }, [selected?.uniqueURI]);

  const subscribe = useCallback((listener) => {
    listeners.current.add(listener);
    return () => listeners.current.delete(listener);
  }, []);

  return (
    <EventsContext.Provider value={{ status, subscribe }}>
      {children}
    </EventsContext.Provider>
  );
}

export function useCampaignEvent(handler) {
  const ctx = useContext(EventsContext);
  if (!ctx) {
    throw new Error("useCampaignEvent must be used within an EventsProvider");
  }
  const { subscribe } = ctx;
  const handlerRef = useRef(handler);
  handlerRef.current = handler;

  useEffect(() => {
    return subscribe((data, event) => handlerRef.current?.(data, event));
  }, [subscribe]);
}
