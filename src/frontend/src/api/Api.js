const API_BASE_URL = import.meta.env.VITE_API_BASE_URL;

export default async function apiRequest(endpoint, options = {}) {
  const headers = options.headers || {};

  // Si no hay un Content-Type y el cuerpo no es un FormData, asignamos application/json
  if (!(options.body instanceof FormData) && !headers["Content-Type"]) {
    headers["Content-Type"] = "application/json";
  }

  const response = await fetch(`${API_BASE_URL}${endpoint}`, {
    mode: 'cors', // Esto indica que la solicitud será una solicitud CORS
    credentials: 'include',
    headers: {
      ...headers,
    },
    ...options,
  });

  // Si la respuesta es un 401 Unauthorized, redirigimos al usuario a la página de inicio de sesión
  if (response.status === 401) {
    window.location.href = "/login";
    //cerrar sesion ususario
    sessionStorage.removeItem("user");
    sessionStorage.removeItem("usuario");
    sessionStorage.removeItem("id");
    sessionStorage.removeItem("role");
    //Elimina el token de las cookies
    document.cookie = `jwt_token=; domain=localhost; path=/;`;
    document.cookie = `role=; domain=localhost; path=/;`;
    return;
  }
    
  // Comprobamos si la respuesta fue correcta (status 200-299)
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

  // Si la respuesta es diferente a 204 No Content, intentamos convertirla en JSON o texto
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
