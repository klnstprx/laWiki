import apiRequest from "./Api.js";


export async function getAllUsers() {
  return apiRequest("/auth");
}

export async function postUser(data) {
  return apiRequest("/auth", {
    method: "POST",
    body: JSON.stringify(data),
  });
}

export async function getUser(id) {
  return apiRequest(`/auth/user?id=${id}`);
}

export async function putUser(id, data) {
  return apiRequest(`/auth/user?id=${id}`, {
    method: "PUT",
    body: JSON.stringify(data),
  });
}

export async function deleteUser(id) {
  return apiRequest(`/auth/users/${encodeURIComponent(id)}`, {
    method: "DELETE",
  });
}