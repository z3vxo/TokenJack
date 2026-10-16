import { getAccess } from "./auth.js";

export async function login(username, password) {
  let res;
  try {
    res = await fetch("/rest/user/login", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ username, password }),
    });
  } catch {
    throw new Error("Cannot reach the server. Check your connection.");
  }

  if (res.ok) {
    return res.json();
  }

  if (res.status === 401) {
    throw new Error("Incorrect username or password.");
  }

  throw new Error("Something went wrong. Please try again.");
}

export async function getCampaigns() {
  let res;
  try {
    res = await fetch("/rest/campaigns/list", {
      headers: { Authorization: `Bearer ${getAccess()}` },
    });
  } catch {
    const err = new Error("Cannot reach the server. Check your connection.");
    err.status = 0;
    throw err;
  }

  if (res.ok) {
    return res.json();
  }

  const err = new Error("Failed to load campaigns.");
  err.status = res.status;
  throw err;
}

export async function getCampaignMeta(id) {
  let res;
  try {
    res = await fetch(`/rest/campaign/meta/${encodeURIComponent(id)}`, {
      headers: { Authorization: `Bearer ${getAccess()}` },
    });
  } catch {
    const err = new Error("Cannot reach the server. Check your connection.");
    err.status = 0;
    throw err;
  }

  if (res.ok) {
    return res.json();
  }

  const err = new Error("Failed to load the campaign overview.");
  err.status = res.status;
  throw err;
}

export async function createCampaign({ name, tenantId, resource, scopes }) {
  let res;
  try {
    res = await fetch("/rest/campaign/new", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${getAccess()}`,
      },
      body: JSON.stringify({
        name,
        tenant_id: tenantId,
        resources: { name: resource, scopes },
      }),
    });
  } catch {
    const err = new Error("Cannot reach the server. Check your connection.");
    err.status = 0;
    throw err;
  }

  if (res.ok) {
    return res.json();
  }

  let message = "Could not create the campaign.";
  try {
    const body = await res.json();
    if (body?.error) message = body.error;
  } catch {
    // ignore non-JSON error bodies
  }

  const err = new Error(message);
  err.status = res.status;
  throw err;
}
