const API_BASE_URL = import.meta.env.VITE_API_BASE_URL;

export default async function apiRequest(endpoint, options = {}) {
  const headers = options.headers || {};

  if (!(options.body instanceof FormData) && !headers["Content-Type"]) {
    headers["Content-Type"] = "application/json";
  }

  const response = await fetch(`${API_BASE_URL}${endpoint}`, {
    headers: {
      ...headers,
    },
    ...options,
  });

  if (!response.ok) {
    let errorMessage = `HTTP error! status: ${response.status}`;
    let errorData;
    try {
      errorData = await response.json();
      if (errorData) {
        errorMessage = errorData.message || errorMessage;
      }
    } catch (e) {
      errorData = null;
    }
    const error = new Error(errorMessage);
    error.status = response.status;
    error.data = errorData;
    throw error;
  }

  if (response.status !== 204) {
    const contentType = response.headers.get("Content-Type");
    if (contentType && contentType.includes("application/json")) {
      return response.json();
    } else {
      return response.text();
    }
  } else {
    return null;
  }
}
