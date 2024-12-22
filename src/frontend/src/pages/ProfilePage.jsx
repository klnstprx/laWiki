import React, { useEffect, useState } from "react";
import { Box, Avatar, Typography, Paper } from "@mui/material";
import { useParams } from "react-router-dom";
import { getUser } from "../api/AuthApi";

const ProfilePage = () => {
  const [user, setUser] = useState(null);
  const { id } = useParams();
  
  //funcion para cargar el usuario con el id que se pasa por parametro
  useEffect(() => {
    console.log(id);
    getUser(id)
      .then((data) => setUser(data))
      .catch(() => setUser(null));  
  } , [id]);

  if (!user) {
    return (
      <Box
        sx={{
          display: "flex",
          justifyContent: "center",
          alignItems: "center",
          height: "100vh",
        }}
      >
        <Typography variant="h5">No has iniciado sesión.</Typography>
      </Box>
    );
  }

  return (
    <Box
      sx={{
        display: "flex",
        justifyContent: "center",
        alignItems: "center",
        height: "100vh",
        backgroundColor: "#f5f5f5",
        padding: 2,
      }}
    >
      <Paper
        elevation={3}
        sx={{
          maxWidth: 400,
          width: "100%",
          padding: 4,
          borderRadius: 2,
        }}
      >
        <Box
          sx={{
            display: "flex",
            flexDirection: "column",
            alignItems: "center",
            gap: 2,
          }}
        >
          {/* Foto de perfil */}
          <Avatar
            alt={user.name}
            src={user.picture}
            sx={{
              width: 100,
              height: 100,
            }}
          />
          {/* Nombre */}
          <Typography variant="h5" fontWeight="bold">
            {user.name}
          </Typography>
          {/* Correo electrónico */}
          <Typography variant="body1" color="text.secondary">
            {user.email}
          </Typography>
        </Box>
      </Paper>
    </Box>
  );
};

export default ProfilePage;
