const API_BASE_URL = import.meta.env.VITE_API_BASE_URL;

export default async function apiRequest(endpoint, options = {}) {
  const headers = options.headers || {};

  // If there's no Content-Type and the body isn't FormData, set it to application/json
  if (!(options.body instanceof FormData) && !headers["Content-Type"]) {
    headers["Content-Type"] = "application/json";
  }

  const response = await fetch(`${API_BASE_URL}${endpoint}`, {
    mode: "cors", // Indicates that the request is a CORS request
    credentials: "include", // Include cookies in the request
    headers: {
      ...headers,
    },
    ...options,
  });

  // Check if the response status is 401 Unauthorized
  if (response.status === 401) {
    const error = new Error("Unauthorized");
    error.status = 401;
    throw error;
  }

  // Check if the response was successful (status 200-299)
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

  // If the response is not 204 No Content, attempt to parse it as JSON or text
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
