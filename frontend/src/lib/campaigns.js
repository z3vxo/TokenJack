const KEY = "tj_campaign";

export function loadCampaign() {
  try {
    const raw = localStorage.getItem(KEY);
    return raw ? JSON.parse(raw) : null;
  } catch {
    return null;
  }
}

export function saveCampaign(campaign) {
  if (campaign) {
    localStorage.setItem(KEY, JSON.stringify(campaign));
  } else {
    localStorage.removeItem(KEY);
  }
}
