import { apiRequest } from "./Api.js";
export async function Login() {
  return apiRequest("/auth/login");
}
