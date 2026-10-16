import {
  createContext,
  useCallback,
  useContext,
  useRef,
  useState,
} from "react";
import { CheckCircle, WarningCircle, X } from "@phosphor-icons/react";

const ToastContext = createContext(null);

const AUTO_DISMISS_MS = 4500;

export function ToastProvider({ children }) {
  const [toasts, setToasts] = useState([]);
  const counter = useRef(0);

  const dismiss = useCallback((id) => {
    setToasts((list) => list.filter((toast) => toast.id !== id));
  }, []);

  const push = useCallback(
    (message, tone) => {
      const id = (counter.current += 1);
      setToasts((list) => [...list, { id, message, tone }]);
      setTimeout(() => dismiss(id), AUTO_DISMISS_MS);
    },
    [dismiss],
  );

  const api = {
    success: (message) => push(message, "success"),
    error: (message) => push(message, "error"),
  };

  return (
    <ToastContext.Provider value={api}>
      {children}
      <div className="pointer-events-none fixed inset-x-0 bottom-4 z-50 flex flex-col items-center gap-2 px-4">
        {toasts.map((toast) => (
          <ToastItem key={toast.id} toast={toast} onDismiss={dismiss} />
        ))}
      </div>
    </ToastContext.Provider>
  );
}

export function useToast() {
  const ctx = useContext(ToastContext);
  if (!ctx) {
    throw new Error("useToast must be used within a ToastProvider");
  }
  return ctx;
}

function ToastItem({ toast, onDismiss }) {
  const success = toast.tone === "success";

  return (
    <div
      role={success ? "status" : "alert"}
      className={`pointer-events-auto flex w-full max-w-sm items-start gap-2.5 rounded-control border px-3.5 py-3 text-sm shadow-menu ${
        success
          ? "border-border bg-card text-foreground"
          : "border-destructive-border bg-destructive-surface text-destructive"
      }`}
    >
      {success ? (
        <CheckCircle
          size={18}
          weight="fill"
          aria-hidden="true"
          className="mt-0.5 shrink-0 text-success"
        />
      ) : (
        <WarningCircle
          size={18}
          weight="fill"
          aria-hidden="true"
          className="mt-0.5 shrink-0"
        />
      )}
      <span className="flex-1 leading-5">{toast.message}</span>
      <button
        type="button"
        onClick={() => onDismiss(toast.id)}
        aria-label="Dismiss notification"
        className="-mr-1 -mt-1 shrink-0 cursor-pointer rounded-control p-1 transition-colors hover:text-foreground focus-visible:ring-[3px] focus-visible:ring-[color:var(--focus-ring-color)]"
      >
        <X size={14} weight="bold" aria-hidden="true" />
      </button>
    </div>
  );
}
