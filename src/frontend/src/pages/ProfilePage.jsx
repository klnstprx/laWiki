import React, { useEffect, useState } from "react";
import { Box, Avatar, Typography, Paper, Rating, TextField, Button, Checkbox } from "@mui/material";
import { useParams } from "react-router-dom";
import { getUser, putUser } from "../api/AuthApi";
import { useToast } from "../context/ToastContext.jsx";

const ProfilePage = () => {
  const [user, setUser] = useState(null);
  const { id } = useParams();
  const { showToast } = useToast();

  const [newRating, setNewRating] = useState(0);
  const [newEnableMails, setNewEnableMails] = useState(false);
  const [mediaRating, setMediaRating] = useState(0);
  const [showRatingForm, setShowRatingForm] = useState(true); // Estado para controlar la visibilidad del formulario
  const isLoggedIn = !!sessionStorage.getItem('user'); // Verifica si el usuario está logueado
  // Obtener el email del usuario logueado desde sessionStorage
  const loggedInUser = JSON.parse(sessionStorage.getItem('user'));
  const loggedInUserEmail = loggedInUser ? loggedInUser.email : null;

  const handleRatingChange = (event, newValue) => {
    setNewRating(newValue);
  };

  const handleEnableMailsChange = (event, newValue) => {
    console.log("onchange del checkbox")
    setNewEnableMails(newValue);
    console.log("New Rating:", newRating);
  }

  const handleSubmit = (event) => {
    event.preventDefault();
    console.log("New Rating:", newRating);

    const updatedValoration = user.valoration ? [...user.valoration, newRating] : [newRating];

    putUser(id, { valoration: updatedValoration })
      .then(() => {
        console.log("Valoración enviada.");
        // Actualizar el estado de la valoración
        setUser((prevUser) => ({ ...prevUser, valoration: updatedValoration }));
        setShowRatingForm(false);
        showToast("Valoración enviada");
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

  const handleSubmitEnableMails = (event) => {
    event.preventDefault();
    console.log("EnableMails:", newEnableMails);

    putUser(id, { enable_mails: newEnableMails })
      .then(() => {
        console.log("Configuración enviada.");
        // Actualizar el estado de las notificaciones por correo
        setUser((prevUser) => ({ ...prevUser, enable_mails: newEnableMails }));
        showToast("Configuración actualizada");
        
      })
      .catch(() => {
        console.error("Error al enviar la configuración.");
      });

  }

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

        // Inicializa el estado `newEnableMails` con el valor actual del usuario
        if (data && typeof data.enable_mails === "boolean") {
          setNewEnableMails(data.enable_mails);
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

        {/* Formulario para activar notificaciones por correo */}
        {isLoggedIn && (
          <Box component="form" onSubmit={handleSubmitEnableMails} sx={{ mt: 4 }}>
            <Typography variant="h6" gutterBottom>
              Activar notificaciones por correo
            </Typography>
            <Box sx={{ display: "flex", alignItems: "center", mb: 2 }}>
              <Typography variant="body1" sx={{ mr: 2 }}>
                Activo: notificaciones por correo
                Inactivo: notificaciones internas
              </Typography>
              <Checkbox
                checked={newEnableMails}
                onChange={handleEnableMailsChange}
                inputProps={{ "aria-label": "Activar notificaciones por correo" }}
              />
            </Box>
            <Button type="submit" variant="contained" fullWidth>
              Aceptar
            </Button>
          </Box>
        )}

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
