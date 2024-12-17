import { useState } from "react";

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
      <h1>Login</h1>
      {error && <p style={{ color: "red" }}>{error}</p>}
      <button onClick={handleLogin}>Login with Google</button>
    </div>
  );
}

export default LoginPage;
