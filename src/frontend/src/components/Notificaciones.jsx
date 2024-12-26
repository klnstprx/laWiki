import React from "react";
import PropTypes from "prop-types";
import {
  Popover,
  List,
  ListItem,
  ListItemText,
  Button,
  Box,
} from "@mui/material";

const Notificaciones = ({
  anchorEl,
  open,
  onClose,
  notifications,
  clearNotifications,
}) => {
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
          {notifications.length > 0 ? (
            notifications.map((notification, index) => (
              <ListItem key={index}>
                <ListItemText primary={notification} />
              </ListItem>
            ))
          ) : (
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
Notificaciones.propTypes = {
    anchorEl: PropTypes.object,
    open: PropTypes.bool.isRequired,
    onClose: PropTypes.func.isRequired,
    notifications: PropTypes.array.isRequired,
    clearNotifications: PropTypes.func.isRequired,
};

export default Notificaciones;
