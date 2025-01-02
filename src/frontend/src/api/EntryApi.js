import { apiRequest, buildQueryString } from "./Api.js";

export async function getAllEntries() {
  return apiRequest("/entries");
}

export async function postEntry(data) {
  return apiRequest("/entries", {
    method: "POST",
    body: JSON.stringify(data),
  });
}

export async function getEntry(id) {
  return apiRequest(`/entries/${encodeURIComponent(id)}`);
}

export async function putEntry(id, data) {
  return apiRequest(`/entries/${encodeURIComponent(id)}`, {
    method: "PUT",
    body: JSON.stringify(data),
  });
}

export async function deleteEntry(id) {
  return apiRequest(`/entries/${encodeURIComponent(id)}`, {
    method: "DELETE",
  });
}

export async function searchEntries(params) {
  const queryString = buildQueryString(params);
  return apiRequest(`/entries/search?${queryString}`);
}

export async function translateEntry(id, targetLang) {
  const queryString = new URLSearchParams({ targetLang }).toString();
  return apiRequest(
    `/entries/${encodeURIComponent(id)}/translate?${queryString}`,
    {
      method: "POST",
    },
  );
}
export async function getUsersByName(name) {
  return apiRequest(`/auth?name=${encodeURIComponent(name)}`);
}
