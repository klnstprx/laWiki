import React, { useState, useEffect } from "react";
import { GoogleLogin } from "@react-oauth/google";
import {jwtDecode} from "jwt-decode";
import {
  Typography,
  Button,
  Box
} from "@mui/material";
import { useNavigate } from "react-router-dom";
import { postUser } from "../api/AuthApi"

const LoginButton = () => {
  const [user, setUser] = useState(null);
  const [addedUser, setAddedUser] = useState(null);
  const navigate = useNavigate();

  useEffect(() => {
    // Verifica si hay un usuario en la sesión al cargar el componente
    const savedUser = sessionStorage.getItem("user");
    if (savedUser) {
      setUser(JSON.parse(savedUser));
    }
  }, []);

  const handleLoginSuccess = async (credentialResponse) => {
    try {
      const decodedUser = jwtDecode(credentialResponse.credential);
      console.log("Usuario autenticado:", decodedUser);

      // Guarda el usuario en la sesión y en el estado
      sessionStorage.setItem("user", JSON.stringify(decodedUser));
      setUser(decodedUser);

      /*
      Se debe añadir usuario a la base de datos con los siguientes datos como mínimo:
      Email         string  
      Name          string
      GivenName     string 
      FamilyName    string 
      Picture       string 
      Role          string 
      Valoration    double
      */

      const user = {
        email: decodedUser.email,
        name: decodedUser.name,
        givenName: decodedUser.given_name,
        familyName: decodedUser.family_name,
        picture: decodedUser.picture,
        role: "user",
        valoration: 0
      };

      const addedUser = await postUser(user);
      setAddedUser(addedUser);
      sessionStorage.setItem("id", addedUser.id);
      console.log('User added:', addedUser);
    
      
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
    navigate("/perfil/" + sessionStorage.getItem("id"));
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