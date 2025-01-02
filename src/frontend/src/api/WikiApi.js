import { apiRequest } from "./Api.js";
export async function getAllWikis() {
  return apiRequest("/wikis");
}

export async function postWiki(data) {
  return apiRequest("/wikis", {
    method: "POST",
    body: JSON.stringify(data),
  });
}

export async function getWiki(id) {
  return apiRequest(`/wikis/${encodeURIComponent(id)}`);
}

export async function putWiki(id, data) {
  return apiRequest(`/wikis/${encodeURIComponent(id)}`, {
    method: "PUT",
    body: JSON.stringify(data),
  });
}

export async function deleteWiki(id) {
  return apiRequest(`/wikis/${encodeURIComponent(id)}`, {
    method: "DELETE",
  });
}

export async function searchWikis(params) {
  const queryString = new URLSearchParams(params).toString();
  return apiRequest(`/wikis/search?${queryString}`);
}

export async function translateWiki(id, targetLang) {
  const queryString = new URLSearchParams({ targetLang }).toString();
  return apiRequest(
    `/wikis/${encodeURIComponent(id)}/translate?${queryString}`,
    {
      method: "POST",
    },
  );
}
