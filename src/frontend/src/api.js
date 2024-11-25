const API_BASE_URL = import.meta.env.VITE_API_BASE_URL;

export async function getAllWikis() {
  const resp = await fetch(`${API_BASE_URL}/wikis`);
  if (!resp.ok) {
    throw new Error("Failed to fetch wikis");
  }
  return resp.json();
}

export async function getWiki(id) {
  const resp = await fetch(`${API_BASE_URL}/wikis/${id}`);
  if (!resp.ok) {
    throw new Error("Failed to fetch wikis");
  }
  return resp.json();
}


export async function getEntrada(id) {
  const resp = await fetch(`${API_BASE_URL}/entries/id?id=${id}`);
  if (!resp.ok) {
    throw new Error("Failed to fetch entries");
  }
  return resp.json();
}

export async function getAllComentariosByVersion(id) {
  const resp = await fetch(`${API_BASE_URL}/comments/version?versionID=${id}`);
  if (!resp.ok) {
    throw new Error("Failed to fetch comments");
  }
  return resp.json();
}

export async function getVersionById(id) {
  const resp = await fetch(`${API_BASE_URL}/versions/id/?id=${id}`);
  if (!resp.ok) {
    throw new Error("Failed to fetch version");
  }
  return resp.json();
}







//other api calls
