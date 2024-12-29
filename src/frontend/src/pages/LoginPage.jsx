import { useState } from "react";
import LoginButton from "../components/LoginButton";

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL;

function LoginPage() {
  const [error, setError] = useState(null);
  const handleLogin = () => {
    try {
      window.location.href = API_BASE_URL + "/auth/login";
    } catch (err) {
      setError("Login failed. Please try again.");
    }
  };

  return (
    <div>
      <h1>Debes haber iniciado sesion.</h1>
    </div>
  );
}

export default LoginPage;
