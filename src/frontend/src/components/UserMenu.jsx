import { useState } from "react";
import { Box, Button, Typography } from "@mui/material";
import { useNavigate } from "react-router-dom";
import NotificationsMenu from "./NotificationsMenu";
import PropTypes from "prop-types";

const UserMenu = ({ user, logout }) => {
  const [anchorEl, setAnchorEl] = useState(null);
  const open = Boolean(anchorEl);
  const navigate = useNavigate();

  const handleLogout = async () => {
    await logout();
    navigate("/");
  };

  const handleNotificationsClick = (event) => {
    setAnchorEl(event.currentTarget);
  };

  const handleNotificationsClose = () => {
    setAnchorEl(null);
  };

  const goToProfile = () => {
    if (user) {
      navigate(`/perfil/${user.id}`);
    }
  };

  return (
    <Box sx={{ display: "flex", alignItems: "center", gap: 2 }}>
      <Button
        variant="contained"
        color="secondary"
        onClick={handleNotificationsClick}
      >
        Notificaciones
      </Button>
      <NotificationsMenu
        anchorEl={anchorEl}
        open={open}
        onClose={handleNotificationsClose}
        user={user}
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

      <Button variant="contained" color="info" onClick={goToProfile}>
        Perfil
      </Button>

      <Button variant="contained" color="error" onClick={handleLogout}>
        Cerrar sesión
      </Button>
    </Box>
  );
};

UserMenu.propTypes = {
  user: PropTypes.object.isRequired,
  logout: PropTypes.func.isRequired,
};

export default UserMenu;
