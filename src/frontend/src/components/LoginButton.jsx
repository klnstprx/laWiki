import React, { useState, useEffect } from "react";
import { GoogleLogin } from "@react-oauth/google";
import {jwtDecode} from "jwt-decode";
import {
  Typography,
  Button,
  Box,
  Popover,
  List,
  ListItem,
  ListItemText,
} from "@mui/material";
import { useNavigate } from "react-router-dom";
import { postUser, getAllUsers } from "../api/AuthApi"

const LoginButton = () => {
  const [user, setUser] = useState(null);
  const [addedUser, setAddedUser] = useState(null);
  const navigate = useNavigate();

  const [anchorEl, setAnchorEl] = useState(null);
  const open = Boolean(anchorEl);
  const notifications = ["Notificación 1", "Notificación 2", "Notificación 3"];

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

      //carga credentialResponse en las cookies con dominio localhost
      document.cookie = `jwt_token=${credentialResponse.credential}; domain=localhost; path=/`;


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
        valoration: []
      };

      {/*}
      const users = await getAllUsers();

      //verifica si el usuario ya existe
      let userExists = false;
      for (let i = 0; i < users.length; i++) {
        if (users[i].email === user.email) {
          console.log("User already registered");
          sessionStorage.setItem("id", users[i].id);
          userExists = true;
          break;
        }
      }

      if (!userExists) {
        console.log("User not registered");
        const addedUser = await postUser(user);
        setAddedUser(addedUser);
        sessionStorage.setItem("id", addedUser.id);
        console.log('User added:', addedUser);
      }

      */}
      //Recarga todas las paginas cuando alguien inicia sesion
      window.location.reload();
      

      
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
    window.location.reload();
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

  return (
    <Box sx={{ display: "flex", alignItems: "center", gap: 2 }}>
      {user ? (
        <>
          <Button
            variant="contained"
            color="info"
            onClick={handleClick}
          >
            Notificaciones
          </Button>
          <Popover
            open={open}
            anchorEl={anchorEl}
            onClose={handleClose}
            anchorOrigin={{
              vertical: 'bottom',
              horizontal: 'left',
            }}
          >
            <List>
              {notifications.map((notification, index) => (
                <ListItem key={index}>
                  <ListItemText primary={notification} />
                </ListItem>
              ))}
            </List>
          </Popover>
          
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
            color="success"
            onClick={goToProfile}
          >
            Perfil
          </Button>
          
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