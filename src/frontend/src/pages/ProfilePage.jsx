import React, { useEffect, useState } from "react";
import { Box, Avatar, Typography, Paper, Rating, TextField, Button } from "@mui/material";
import { useParams } from "react-router-dom";
import { getUser, putUser } from "../api/AuthApi";

const ProfilePage = () => {
  const [user, setUser] = useState(null);
  const { id } = useParams();

  const [newRating, setNewRating] = useState(0);
  const [mediaRating, setMediaRating] = useState(0);
  const [showRatingForm, setShowRatingForm] = useState(true); // Estado para controlar la visibilidad del formulario
  const isLoggedIn = !!sessionStorage.getItem('user'); // Verifica si el usuario está logueado
  // Obtener el email del usuario logueado desde sessionStorage
  const loggedInUser = JSON.parse(sessionStorage.getItem('user'));
  const loggedInUserEmail = loggedInUser ? loggedInUser.email : null;

  const handleRatingChange = (event, newValue) => {
    setNewRating(newValue);
  };

  const handleSubmit = (event) => {
    event.preventDefault();
    console.log("New Rating:", newRating);
    // Aquí podrías enviar los datos a tu API
    const updatedValoration = user.valoration ? [...user.valoration, newRating] : [newRating];

    putUser(id, { valoration: updatedValoration })
      .then(() => {
        console.log("Valoración enviada.");
        // Actualizar el estado de la valoración
        setUser((prevUser) => ({ ...prevUser, valoration: updatedValoration }));
        setShowRatingForm(false);
      })
      .catch(() => {
        console.error("Error al enviar la valoración.");
      });

    setNewRating(0);
      //actualiza el campo setMediaRating
      const sum = updatedValoration.reduce((acc, val) => acc + val, 0);
      const media = sum / updatedValoration.length;
      setMediaRating(media);
      

  };

  useEffect(() => {
    console.log(id);
    getUser(id)
      .then((data) => {
        setUser(data);
        if (data && data.valoration && data.valoration.length > 0) {
          const sum = data.valoration.reduce((acc, val) => acc + val, 0);
          const media = sum / data.valoration.length;
          setMediaRating(media);
        }
      })
      .catch(() => setUser(null));
  }, [id]);

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
          {/* Valoración */}
          <Rating
            name="user-rating"
            value={mediaRating}
            precision={0.5}
            readOnly
          />
          ({user.valoration.length})
        </Box>
        

        {/* Formulario para valorar */}
        {isLoggedIn && showRatingForm && user.email != loggedInUserEmail && (
        <Box component="form" onSubmit={handleSubmit} sx={{ mt: 4 }}>
          <Typography variant="h6" gutterBottom>
            Valorar al usuario
          </Typography>
          <Rating
            name="new-rating"
            value={newRating}
            onChange={handleRatingChange}
            precision={0.5}
            sx={{ mb: 2 }}
          />
          <Button type="submit" variant="contained" fullWidth>
            Enviar valoración
          </Button>
        </Box>
        )}
      </Paper>
    </Box>
  );
};

export default ProfilePage;
