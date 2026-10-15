const ACCESS_KEY = "tj_access";
const REFRESH_KEY = "tj_refresh";

export function getAccess() {
  return localStorage.getItem(ACCESS_KEY);
}

export function isAuthed() {
  return Boolean(getAccess());
}

export function setTokens({ access, refresh }) {
  if (access) localStorage.setItem(ACCESS_KEY, access);
  if (refresh) localStorage.setItem(REFRESH_KEY, refresh);
}

export function clearTokens() {
  localStorage.removeItem(ACCESS_KEY);
  localStorage.removeItem(REFRESH_KEY);
}
