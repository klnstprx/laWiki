import { Button } from "@mui/material";

const LoginButton = () => {
  const handleLogin = () => {
    window.location.href = `${import.meta.env.VITE_API_BASE_URL}/auth/login`;
  };

  return (
    <Button variant="contained" color="primary" onClick={handleLogin}>
      Iniciar sesión
    </Button>
  );
};

export default LoginButton;
