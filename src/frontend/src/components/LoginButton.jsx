import React, { useState, useEffect } from "react";
import { GoogleLogin } from "@react-oauth/google";
import { jwtDecode } from "jwt-decode";
import {
  Typography,
  Button,
  Box,
} from "@mui/material";
import { useNavigate } from "react-router-dom";
import { postUser, getAllUsers } from "../api/AuthApi";
import Notificaciones from "./Notificaciones"; // Importa el nuevo componente
import { putUser } from "../api/AuthApi";

const LoginButton = () => {
  const [user, setUser] = useState(null); // Estado para el usuario autenticado
  const [usuario, setUsuario] = useState({ notifications: [] }); // Estado para el usuario con notificaciones
  const navigate = useNavigate();

  const [anchorEl, setAnchorEl] = useState(null);
  const open = Boolean(anchorEl);

  // Carga inicial del usuario y sus notificaciones desde sessionStorage
  useEffect(() => {
    const savedUser = sessionStorage.getItem("user");
    const savedUsuario = sessionStorage.getItem("usuario");

    if (savedUser) {
      setUser(JSON.parse(savedUser));
    }
    if (savedUsuario) {
      setUsuario(JSON.parse(savedUsuario));
    }
  }, []);

  const handleLoginSuccess = async (credentialResponse) => {
    try {
      const decodedUser = jwtDecode(credentialResponse.credential);
      console.log("Usuario autenticado:", decodedUser);

      // Guarda el usuario en la sesión y en el estado
      sessionStorage.setItem("user", JSON.stringify(decodedUser));
      setUser(decodedUser);

      //carga credentialResponse en las cookies con dominio localhost
      document.cookie = `jwt_token=${credentialResponse.credential}; domain=localhost; path=/`;

      const user = {
        email: decodedUser.email,
        name: decodedUser.name,
        given_name: decodedUser.given_name,
        family_name: decodedUser.family_name,
        picture: decodedUser.picture,
        locale: decodedUser.locale,
        email_verified: decodedUser.email_verified,
        role: "redactor",
        valoration: [],
        notifications: [],
        enable_emails: false,
      };

      const users = await getAllUsers();

      if (users != null) {
      
      // Verifica si el usuario ya existe
      let userExists = false;
      for (let i = 0; i < users.length; i++) {
        if (users[i].email === user.email) {
          console.log("User already registered");
          sessionStorage.setItem("id", users[i].id);
          setUsuario(users[i]); // Actualiza el estado con el usuario existente
          sessionStorage.setItem("usuario", JSON.stringify(users[i]));
          document.cookie = `role=${users[i].role}; domain=localhost; path=/`;
          userExists = true;
          break;
        }
      }

      if (!userExists) {
        console.log("User not registered");
        const addedUser = await postUser(user);
        sessionStorage.setItem("id", addedUser.id);
        setUsuario(addedUser); // Actualiza el estado con el nuevo usuario
        sessionStorage.setItem("usuario", JSON.stringify(addedUser));
        document.cookie = `role=${addedUser.role}; domain=localhost; path=/`;
      }


    } else {
      console.log("User not registered");
      const addedUser = await postUser(user);
      sessionStorage.setItem("id", addedUser.id);
      setUsuario(addedUser); // Actualiza el estado con el nuevo usuario
      sessionStorage.setItem("usuario", JSON.stringify(addedUser));
      document.cookie = `role=${addedUser.role}; domain=localhost; path=/`;
    }

    //si la direccion actual es /login, recarga la pagina
    if (window.location.pathname === "/login") {
      window.location.href = "/";
    } else {
      window.location.reload();
    }


    } catch (error) {
      console.error("Error al procesar las credenciales:", error);
    }
  };

  const handleLoginError = () => {
    console.error("Error al iniciar sesión");
    alert("Hubo un problema al iniciar sesión. Inténtalo de nuevo.");
  };

  const handleLogout = () => {
    sessionStorage.removeItem("user");
    sessionStorage.removeItem("usuario");
    sessionStorage.removeItem("id");
    //Elimina el token de las cookies
    document.cookie = `jwt_token=; domain=localhost; path=/; expires=Thu, 01 Jan 1970 00:00:00 UTC`;
    setUser(null);
    setUsuario({ notifications: [] });
  };

  const goToProfile = () => {
    navigate("/perfil/" + sessionStorage.getItem("id"));
  };

  const handleClick = (event) => {
    setAnchorEl(event.currentTarget);
  };

  const handleClose = () => {
    setAnchorEl(null);
  };

  const clearNotifications = () => {
    const updatedUsuario = { ...usuario, notifications: [] };
    setUsuario(updatedUsuario);
    sessionStorage.setItem("usuario", JSON.stringify(updatedUsuario));

    putUser(sessionStorage.getItem("id"), { notifications: [] })
          .then(() => {
            console.log("Notificaciones eliminadas.");
            setUser((prevUser) => ({ ...prevUser, notifications: [] }));
            
          })
          .catch(() => {
            console.error("Error al enviar la configuración.");
          });
  };

  return (
    <Box sx={{ display: "flex", alignItems: "center", gap: 2 }}>
      {user ? (
        <>
          <Button variant="contained" color="info" onClick={handleClick}>
            Notificaciones
          </Button>
          <Notificaciones
            anchorEl={anchorEl}
            open={open}
            onClose={handleClose}
            notifications={usuario.notifications}
            clearNotifications={clearNotifications}
          />

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

          <Button variant="contained" color="success" onClick={goToProfile}>
            Perfil
          </Button>

          <Button variant="contained" color="warning" onClick={handleLogout}>
            Cerrar sesión
          </Button>
        </>
      ) : (
        <GoogleLogin onSuccess={handleLoginSuccess} onError={handleLoginError} />
      )}
    </Box>
  );
};

export default LoginButton;
