const API_BASE_URL = import.meta.env.VITE_API_BASE_URL;

export default async function apiRequest(endpoint, options = {}) {
  const resp = await fetch(`${API_BASE_URL}${endpoint}`, {
    headers: {
      "Content-Type": "application/json",
      ...options.headers,
    },
    ...options,
  });

  if (!resp.ok) {
    let errorMessage = `HTTP error! status: ${resp.status}`;
    const errorData = await resp.json().catch(() => null);
    if (errorData) {
      errorMessage = errorData.message || errorMessage;
    }
    const error = new Error(errorMessage);
    error.status = resp.status;
    error.data = errorData;
    throw error;
  }

  // delete no tienen contenido
  if (resp.status !== 204) {
    return resp.json();
  } else {
    return null;
  }
}
