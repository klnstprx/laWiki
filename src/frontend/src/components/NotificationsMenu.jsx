import { useState } from "react";
import PropTypes from "prop-types";
import {
  Box,
  Button,
  List,
  ListItem,
  ListItemText,
  Popover,
} from "@mui/material";
import { putUser } from "../api/AuthApi";

const NotificationsMenu = ({ anchorEl, open, onClose, user }) => {
  const [notifications, setNotifications] = useState(user.notifications || []);

  const clearNotifications = async () => {
    // Update the user's notifications to an empty array in the backend
    try {
      await putUser(user.id, { notifications: [] });
      setNotifications([]);
      console.log("Notifications cleared in the backend.");
    } catch (error) {
      console.error("Error clearing notifications in the backend:", error);
    }
  };

  return (
    <Popover
      open={open}
      anchorEl={anchorEl}
      onClose={onClose}
      anchorOrigin={{
        vertical: "bottom",
        horizontal: "left",
      }}
    >
      <Box sx={{ padding: 2, minWidth: 300 }}>
        <List>
          {notifications.length > 0
            ? (
              notifications.map((notification, index) => (
                <ListItem key={index}>
                  <ListItemText primary={notification} />
                </ListItem>
              ))
            )
            : (
              <ListItem>
                <ListItemText primary="No hay notificaciones" />
              </ListItem>
            )}
        </List>
        {notifications.length > 0 && (
          <Button
            variant="contained"
            color="error"
            fullWidth
            onClick={clearNotifications}
          >
            Eliminar todas las notificaciones
          </Button>
        )}
      </Box>
    </Popover>
  );
};

NotificationsMenu.propTypes = {
  anchorEl: PropTypes.object,
  open: PropTypes.bool.isRequired,
  onClose: PropTypes.func.isRequired,
  user: PropTypes.object.isRequired,
};

export default NotificationsMenu;
