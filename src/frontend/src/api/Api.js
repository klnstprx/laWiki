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
    try {
      const errorData = await resp.json();
      errorMessage = errorData.message || errorMessage;
    } catch {
      // Ignore JSON parsing errors
    }
    throw new Error(errorMessage);
  }

  // delete no tienen contenido
  if (resp.status !== 204) {
    return resp.json();
  } else {
    return null;
  }
}
