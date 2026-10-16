export function campaignEventsUrl(uniqueURI, token) {
  const params = new URLSearchParams({ auth: token });
  return `/rest/events/${encodeURIComponent(uniqueURI)}?${params.toString()}`;
}
