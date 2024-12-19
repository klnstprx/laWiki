import React, { useState, useEffect } from "react";
import { GoogleLogin } from "@react-oauth/google";
import {jwtDecode} from "jwt-decode";
import {
  Typography,
  Button,
  Box
} from "@mui/material";
import { useNavigate } from "react-router-dom";

const LoginButton = () => {
  const [user, setUser] = useState(null);
  const navigate = useNavigate();

  useEffect(() => {
    // Verifica si hay un usuario en la sesión al cargar el componente
    const savedUser = sessionStorage.getItem("user");
    if (savedUser) {
      setUser(JSON.parse(savedUser));
    }
  }, []);

  const handleLoginSuccess = (credentialResponse) => {
    try {
      const decodedUser = jwtDecode(credentialResponse.credential);
      console.log("Usuario autenticado:", decodedUser);

      // Guarda el usuario en la sesión y en el estado
      sessionStorage.setItem("user", JSON.stringify(decodedUser));
      setUser(decodedUser);

      // Aquí se debe añadir el usuario a la base de datos

    } catch (error) {
      console.error("Error al procesar las credenciales:", error);
    }
  };

  const handleLoginError = () => {
    console.error("Error al iniciar sesión");
    alert("Hubo un problema al iniciar sesión. Inténtalo de nuevo.");
  };

  const handleLogout = () => {
    // Elimina el usuario de la sesión y actualiza el estado
    sessionStorage.removeItem("user");
    setUser(null);
  };

  const goToProfile = () => {
    navigate("/perfil");
  };

  return (
    <Box sx={{ display: "flex", alignItems: "center", gap: 2 }}>
      {user ? (
        <>
          <Button
            variant="contained"
            color="success"
            onClick={goToProfile}
          >
            Perfil
          </Button>
          <Typography
            variant="body1"
            noWrap
            sx={{
              whiteSpace: "nowrap",
              overflow: "hidden",
              textOverflow: "ellipsis",
            }}
          >
            Bienvenido, {user.name}
          </Typography>
          <Button
            variant="contained"
            color="warning"
            onClick={handleLogout}
          >
            Cerrar sesión
          </Button>
        </>
      ) : (
        <GoogleLogin
          onSuccess={handleLoginSuccess}
          onError={handleLoginError}
        />
      )}
    </Box>
  );
};

export default LoginButton;