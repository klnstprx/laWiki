import { apiRequest } from "./Api.js";

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

export async function getUserByEmail(email) {
  try {
    return await apiRequest(`/auth/user/email?email=${email}`);
  } catch (error) {
    if (error.status === 404) {
      return null;
    }
    throw error;
  }
}

export async function getUsersByIds(ids) {
  const idsParam = ids.map(encodeURIComponent).join(",");
  return apiRequest(`/auth/user/ids?ids=${idsParam}`);
}

export async function getUsersByName(name) {
  return await apiRequest(`/auth?name=${encodeURIComponent(name)}`);
}

export async function putUser(id, data) {
  return apiRequest(`/auth/user?id=${id}`, {
    method: "PUT",
    body: JSON.stringify(data),
  });
}

export async function deleteUser(id) {
  return apiRequest(`/auth/user/${encodeURIComponent(id)}`, {
    method: "DELETE",
  });
}
