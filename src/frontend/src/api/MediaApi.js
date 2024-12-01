import apiRequest from "./Api.js";

export async function getAllMedia() {
  return apiRequest("/media");
}

export async function postMedia(data) {
  return apiRequest("/media", {
    method: "POST",
    body: data,
  });
}

export async function getMedia(id) {
  return apiRequest(`/media/${encodeURIComponent(id)}`);
}

export async function putMedia(id, data) {
  return apiRequest(`/media/${encodeURIComponent(id)}`, {
    method: "PUT",
    body: data,
  });
}

export async function deleteMedia(id) {
  return apiRequest(`/media/${encodeURIComponent(id)}`, {
    method: "DELETE",
  });
}

export async function searchMedia(params) {
  const queryString = new URLSearchParams(params).toString();
  return apiRequest(`/media/search?${queryString}`);
}
