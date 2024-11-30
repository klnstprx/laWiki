import apiRequest from "./Api.js";

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
  const queryString = new URLSearchParams(params).toString();
  return apiRequest(`/entries/search?${queryString}`);
}
