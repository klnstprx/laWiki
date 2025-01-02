import { useEffect, useState } from "react";
import { GoogleLogin } from "@react-oauth/google";
import { jwtDecode } from "jwt-decode";
import {
  Avatar,
  Box,
  Button,
  IconButton,
  Typography,
  useMediaQuery,
  useTheme,
} from "@mui/material";
import { useNavigate } from "react-router-dom";
import { getUserByEmail, postUser, putUser } from "../api/AuthApi";
import Notificaciones from "./Notificaciones";
import {
  Logout as LogoutIcon,
  Notifications as NotificationsIcon,
} from "@mui/icons-material";

const LoginButton = () => {
  const [googleUser, setGoogleUser] = useState(null); // State for authenticated Google user
  const [appUser, setAppUser] = useState({ notifications: [] }); // State for user in our app
  const navigate = useNavigate();

  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down("md"));
  const [anchorEl, setAnchorEl] = useState(null);
  const open = Boolean(anchorEl);

  // Initial load of user and notifications from sessionStorage
  useEffect(() => {
    const savedGoogleUser = sessionStorage.getItem("googleUser");
    const savedAppUser = sessionStorage.getItem("appUser");

    if (savedGoogleUser) {
      setGoogleUser(JSON.parse(savedGoogleUser));
    }
    if (savedAppUser) {
      setAppUser(JSON.parse(savedAppUser));
    }
  }, []);

  const handleLoginSuccess = async (credentialResponse) => {
    try {
      const decodedGoogleUser = jwtDecode(credentialResponse.credential);
      console.log("Authenticated user:", decodedGoogleUser);

      // Save the Google user in session and state
      sessionStorage.setItem("googleUser", JSON.stringify(decodedGoogleUser));
      setGoogleUser(decodedGoogleUser);

      // Store the JWT token in cookies
      document.cookie = `jwt_token=${credentialResponse.credential}; path=/`;

      const newUser = {
        email: decodedGoogleUser.email,
        name: decodedGoogleUser.name,
        given_name: decodedGoogleUser.given_name,
        family_name: decodedGoogleUser.family_name,
        picture: decodedGoogleUser.picture,
        locale: decodedGoogleUser.locale,
        email_verified: decodedGoogleUser.email_verified,
        role: "redactor",
        valoration: [],
        notifications: [],
        enable_emails: false,
      };

      // Check if the user already exists
      const existingUser = await getUserByEmail(newUser.email);

      if (existingUser) {
        console.log("User already registered");
        setAppUser(existingUser);
        sessionStorage.setItem("appUser", JSON.stringify(existingUser));
      } else {
        console.log("User not registered");
        const addedUser = await postUser(newUser);
        setAppUser(addedUser);
        sessionStorage.setItem("appUser", JSON.stringify(addedUser));
      }
      window.location.reload();
    } catch (error) {
      console.error("Error processing credentials:", error);
    }
  };

  const handleLoginError = () => {
    console.error("Error logging in");
    alert("There was a problem logging in. Please try again.");
  };

  const handleLogout = () => {
    sessionStorage.removeItem("googleUser");
    sessionStorage.removeItem("appUser");
    // Remove tokens from cookies
    document.cookie = `jwt_token=; path=/;`;
    setGoogleUser(null);
    setAppUser({ notifications: [] });
    window.location.reload();
  };

  const goToProfile = () => {
    navigate(`/perfil/${appUser.id}`);
  };

  const handleClick = (event) => {
    setAnchorEl(event.currentTarget);
  };

  const handleClose = () => {
    setAnchorEl(null);
  };

  const clearNotifications = () => {
    const updatedAppUser = { ...appUser, notifications: [] };

    // Update state and sessionStorage immediately for better user experience
    setAppUser(updatedAppUser);
    sessionStorage.setItem("appUser", JSON.stringify(updatedAppUser));

    // Call the API to update in the database
    putUser(appUser.id, { notifications: [] })
      .then(() => {
        console.log("Notifications cleared in the database.");
      })
      .catch((error) => {
        console.error("Error clearing notifications in the database:", error);
        // Optionally: restore local notifications if an error occurs
        setAppUser(appUser);
        sessionStorage.setItem("appUser", JSON.stringify(appUser));
      });
  };

  return (
    <Box sx={{ display: "flex", alignItems: "center", gap: 2 }}>
      {googleUser
        ? (
          <>
            {isMobile
              ? (
                <Button
                  variant="contained"
                  color="secondary"
                  onClick={handleClick}
                  sx={{
                    minWidth: "auto",
                    p: 1,
                  }}
                >
                  <NotificationsIcon
                    sx={{ display: { xs: "flex", md: "none" } }}
                  />
                </Button>
              )
              : (
                <Button
                  variant="contained"
                  color="secondary"
                  onClick={handleClick}
                  startIcon={
                    <NotificationsIcon
                      sx={{ display: { xs: "none", md: "inline" } }}
                    />
                  }
                  sx={{
                    minWidth: "auto",
                    p: 1,
                  }}
                >
                  <Typography
                    noWrap
                    sx={{ display: { xs: "none", md: "inline" } }}
                  >
                    Notificaciones
                  </Typography>
                </Button>
              )}
            <Notificaciones
              anchorEl={anchorEl}
              open={open}
              onClose={handleClose}
              notifications={appUser.notifications || []}
              clearNotifications={clearNotifications}
            />

            <IconButton
              onClick={goToProfile}
              sx={{
                p: 1,
                minWidth: "auto",
              }}
            >
              <Avatar
                sx={{
                  backgroundColor: "primary.main",
                  "&:hover": {
                    backgroundColor: "primary.dark",
                  },
                  boxShadow: (theme) => theme.shadows[2],
                }}
                alt={googleUser.name}
                src={googleUser.picture}
              />
            </IconButton>

            {isMobile
              ? (
                <Button
                  variant="contained"
                  color="error"
                  onClick={handleLogout}
                  sx={{
                    p: 1,
                    minWidth: "auto",
                  }}
                >
                  <LogoutIcon sx={{ display: { xs: "flex", md: "none" } }} />
                </Button>
              )
              : (
                <Button
                  variant="contained"
                  color="error"
                  onClick={handleLogout}
                  sx={{
                    p: 1,
                    minWidth: "auto",
                  }}
                  startIcon={
                    <LogoutIcon
                      sx={{ display: { xs: "none", md: "inline" } }}
                    />
                  }
                >
                  <Typography
                    noWrap
                    sx={{ display: { xs: "none", md: "inline" } }}
                  >
                    Cerrar sesión
                  </Typography>
                </Button>
              )}
          </>
        )
        : (
          <GoogleLogin
            onSuccess={handleLoginSuccess}
            onError={handleLoginError}
          />
        )}
    </Box>
  );
};

export default LoginButton;
