import { Routes, Route, Navigate } from "react-router-dom";
import Login from "./pages/Login.jsx";
import Overview from "./pages/Overview.jsx";
import Campaigns from "./pages/Campaigns.jsx";
import Placeholder from "./pages/Placeholder.jsx";
import AppLayout from "./components/AppLayout.jsx";
import { ToastProvider } from "./components/Toast.jsx";
import { isAuthed } from "./lib/auth.js";
import { CampaignProvider } from "./context/CampaignContext.jsx";
import { EventsProvider } from "./context/EventsContext.jsx";

function RequireAuth({ children }) {
  return isAuthed() ? children : <Navigate to="/login" replace />;
}

export default function App() {
  return (
    <ToastProvider>
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route
          path="/"
          element={
            <RequireAuth>
              <CampaignProvider>
                <EventsProvider>
                  <AppLayout />
                </EventsProvider>
              </CampaignProvider>
            </RequireAuth>
          }
        >
          <Route index element={<Navigate to="/overview" replace />} />
          <Route path="overview" element={<Overview />} />
          <Route path="campaigns" element={<Campaigns />} />
          <Route path="automation" element={<Placeholder title="Automation" />} />
          <Route path="tokens" element={<Placeholder title="Tokens" />} />
          <Route path="settings" element={<Placeholder title="Settings" />} />
        </Route>
        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </ToastProvider>
  );
}
