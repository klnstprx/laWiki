import { apiRequest, buildQueryString } from "./Api.js";

export async function getAllComments() {
  return apiRequest("/comments");
}

export async function postComment(data) {
  return apiRequest("/comments", {
    method: "POST",
    body: JSON.stringify(data),
  });
}

export async function getComment(id) {
  return apiRequest(`/comments/${encodeURIComponent(id)}`);
}

export async function putComment(id, data) {
  return apiRequest(`/comments/${encodeURIComponent(id)}`, {
    method: "PUT",
    body: JSON.stringify(data),
  });
}

export async function deleteComment(id) {
  return apiRequest(`/comments/${encodeURIComponent(id)}`, {
    method: "DELETE",
  });
}

export async function searchComments(params) {
  const queryString = buildQueryString(params);
  return apiRequest(`/comments/search?${queryString}`);
}
