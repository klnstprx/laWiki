import { useEffect, useState } from "react";
import {
  Avatar,
  Box,
  Button,
  Checkbox,
  FormControl,
  InputLabel,
  MenuItem,
  Paper,
  Rating,
  Select,
  Typography,
} from "@mui/material";
import { useParams } from "react-router-dom";
import { getUser, putUser } from "../api/AuthApi";
import { useToast } from "../context/ToastContext.jsx";
import { useAuth } from "../context/AuthContext";

const ProfilePage = () => {
  const [userProfile, setUserProfile] = useState(null); // Renamed to avoid confusion with logged-in user
  const { id } = useParams();
  const { showToast } = useToast();

  const [newRating, setNewRating] = useState(0);
  const [newEnableMails, setNewEnableMails] = useState(false);
  const [mediaRating, setMediaRating] = useState(0);
  const [showRatingForm, setShowRatingForm] = useState(true);
  const [selectedRole, setSelectedRole] = useState("");

  const { user: loggedInUser, setUser: setLoggedInUser } = useAuth(); // Use setUser to update AuthContext
  const isLoggedIn = !!loggedInUser;
  const loggedInUserEmail = loggedInUser?.email || null;

  const handleRatingChange = (event, newValue) => {
    setNewRating(newValue);
  };

  const handleEnableMailsChange = (event) => {
    setNewEnableMails(event.target.checked);
  };

  const handleRoleChange = (event) => {
    setSelectedRole(event.target.value);
  };

  const handleSubmit = (event) => {
    event.preventDefault();

    const updatedValoration = userProfile.valoration
      ? [...userProfile.valoration, newRating]
      : [newRating];

    putUser(id, { valoration: updatedValoration })
      .then(() => {
        // Update the user profile state
        setUserProfile((prevUser) => ({
          ...prevUser,
          valoration: updatedValoration,
        }));
        setShowRatingForm(false);
        showToast("Tu valoración ha sido enviada correctamente.", "success");

        // Recalculate media rating
        const sum = updatedValoration.reduce((acc, val) => acc + val, 0);
        const media = sum / updatedValoration.length;
        setMediaRating(media);
      })
      .catch(() => {
        console.error("Error al enviar la valoración.");
        showToast("Hubo un error al enviar tu valoración.", "error");
      });

    setNewRating(0);
  };

  const handleSubmitEnableMails = (event) => {
    event.preventDefault();

    putUser(id, { enable_mails: newEnableMails })
      .then(() => {
        // Update the user profile state
        setUserProfile((prevUser) => ({
          ...prevUser,
          enable_mails: newEnableMails,
        }));
        showToast(
          "La configuración de notificaciones ha sido actualizada correctamente.",
          "success",
        );
      })
      .catch(() => {
        console.error("Error al enviar la configuración.");
        showToast(
          "Hubo un error al actualizar la configuración de notificaciones.",
          "error",
        );
      });
  };

  const handleSubmitRoleChange = (event) => {
    event.preventDefault();

    putUser(id, { role: selectedRole })
      .then(() => {
        // Update the user profile state
        setUserProfile((prevUser) => ({ ...prevUser, role: selectedRole }));

        // If the logged-in user is changing their own role, update the AuthContext
        if (id === loggedInUser.id) {
          setLoggedInUser((prevUser) => ({ ...prevUser, role: selectedRole }));
        }

        showToast(
          "El rol del usuario ha sido actualizado correctamente.",
          "success",
        );
      })
      .catch(() => {
        console.error("Error al actualizar el rol.");
        showToast("Hubo un error al actualizar el rol del usuario.", "error");
      });
  };

  useEffect(() => {
    getUser(id)
      .then((data) => {
        setUserProfile(data);
        if (data && data.valoration && data.valoration.length > 0) {
          const sum = data.valoration.reduce((acc, val) => acc + val, 0);
          const media = sum / data.valoration.length;
          setMediaRating(media);
        }

        // Initialize `newEnableMails` with the user's current setting
        if (data && typeof data.enable_mails === "boolean") {
          setNewEnableMails(data.enable_mails);
        }

        // Initialize `selectedRole` with the user's current role
        if (data && data.role) {
          setSelectedRole(data.role);
        }
      })
      .catch(() => setUserProfile(null));
  }, [id]);

  if (!userProfile) {
    return (
      <Box
        sx={{
          display: "flex",
          justifyContent: "center",
          alignItems: "center",
          height: "100vh",
        }}
      >
        <Typography variant="h5">Usuario no encontrado.</Typography>
      </Box>
    );
  }

  return (
    <Box
      sx={{
        display: "flex",
        justifyContent: "center",
        alignItems: "center",
        minHeight: "100vh", // Changed to minHeight to avoid overflow issues
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
          {/* Profile Picture */}
          <Avatar
            alt={userProfile.name}
            src={userProfile.picture}
            sx={{
              width: 100,
              height: 100,
            }}
          />
          {/* Name */}
          <Typography variant="h5" fontWeight="bold">
            <span>
              {userProfile.name} ({userProfile.valoration.length})
            </span>
          </Typography>
          {/* Email */}
          <Typography variant="body1" color="text.secondary">
            {userProfile.email}
          </Typography>
          {/* Rating */}
          <Rating
            name="user-rating"
            value={mediaRating}
            precision={0.5}
            readOnly
          />
        </Box>

        {/* Enable Emails Form */}
        {isLoggedIn && userProfile.email === loggedInUserEmail && (
          <Box
            component="form"
            onSubmit={handleSubmitEnableMails}
            sx={{ mt: 4 }}
          >
            <Typography variant="h6" gutterBottom>
              Activar notificaciones por correo
            </Typography>
            <Box sx={{ display: "flex", alignItems: "center", mb: 2 }}>
              <Typography variant="body1" sx={{ mr: 2 }}>
                Activado: notificaciones por correo. Desactivado: notificaciones
                internas.
              </Typography>
              <Checkbox
                checked={newEnableMails}
                onChange={handleEnableMailsChange}
                inputProps={{
                  "aria-label": "Activar notificaciones por correo",
                }}
              />
            </Box>
            <Button type="submit" variant="contained" fullWidth>
              Aceptar
            </Button>
          </Box>
        )}

        {/* Rating Form */}
        {isLoggedIn &&
          showRatingForm &&
          userProfile.email !== loggedInUserEmail && (
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

        {/* Role Change Form */}
        {isLoggedIn && loggedInUser.role === "admin" && (
          <Box
            component="form"
            onSubmit={handleSubmitRoleChange}
            sx={{ mt: 4 }}
          >
            <Typography variant="h6" gutterBottom>
              Cambiar Rol del Usuario
            </Typography>
            <FormControl fullWidth sx={{ mb: 2 }}>
              <InputLabel id="role-select-label">Rol</InputLabel>
              <Select
                labelId="role-select-label"
                id="role-select"
                value={selectedRole}
                label="Rol"
                onChange={handleRoleChange}
              >
                <MenuItem value="admin">admin</MenuItem>
                <MenuItem value="editor">editor</MenuItem>
                <MenuItem value="redactor">redactor</MenuItem>
              </Select>
            </FormControl>
            <Button type="submit" variant="contained" fullWidth>
              Cambiar rol
            </Button>
          </Box>
        )}
      </Paper>
    </Box>
  );
};

export default ProfilePage;
