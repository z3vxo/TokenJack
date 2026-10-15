import { Routes, Route, Navigate } from "react-router-dom";
import Login from "./pages/Login.jsx";
import Placeholder from "./pages/Placeholder.jsx";
import AppLayout from "./components/AppLayout.jsx";
import { isAuthed } from "./lib/auth.js";
import { CampaignProvider } from "./context/CampaignContext.jsx";

function RequireAuth({ children }) {
  return isAuthed() ? children : <Navigate to="/login" replace />;
}

export default function App() {
  return (
    <Routes>
      <Route path="/login" element={<Login />} />
      <Route
        path="/"
        element={
          <RequireAuth>
            <CampaignProvider>
              <AppLayout />
            </CampaignProvider>
          </RequireAuth>
        }
      >
        <Route index element={<Navigate to="/campaigns" replace />} />
        <Route path="campaigns" element={<Placeholder title="Campaigns" />} />
        <Route path="automation" element={<Placeholder title="Automation" />} />
        <Route path="codes" element={<Placeholder title="Codes" />} />
        <Route path="settings" element={<Placeholder title="Settings" />} />
      </Route>
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  );
}
