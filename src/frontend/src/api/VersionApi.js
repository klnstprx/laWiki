import apiRequest from "./Api.js";

export async function getAllVersions() {
  return apiRequest("/versions");
}

export async function postVersion(data) {
  return apiRequest("/versions", {
    method: "POST",
    body: JSON.stringify(data),
  });
}

export async function getVersion(id) {
  return apiRequest(`/versions/${encodeURIComponent(id)}`);
}

export async function putVersion(id, data) {
  return apiRequest(`/versions/${encodeURIComponent(id)}`, {
    method: "PUT",
    body: JSON.stringify(data),
  });
}

export async function deleteVersion(id) {
  return apiRequest(`/versions/${encodeURIComponent(id)}`, {
    method: "DELETE",
  });
}

export async function searchVersions(params) {
  const queryString = new URLSearchParams(params).toString();
  return apiRequest(`/versions/search?${queryString}`);
}

export async function translateVersion(id, targetLang) {
  const queryString = new URLSearchParams({ targetLang }).toString();
  return apiRequest(`/versions/${encodeURIComponent(id)}/translate?${queryString}`, {
    method: "POST",
  });
}