import { Link } from "react-router-dom";
import PropTypes from "prop-types";
import {
  Card,
  CardActionArea,
  CardContent,
  IconButton,
  Typography,
} from "@mui/material";
import Grid from "@mui/joy/Grid";
import DeleteIcon from "@mui/icons-material/Delete";
import ConfirmationModal from "../components/ConfirmationModal.jsx";
import { useState } from "react";

import { useAuth } from "../context/AuthContext";

const EntradaCard = ({
  id,
  title,
  author,
  createdAt,
  onEntradaClick,
  onDelete,
}) => {
  const { user } = useAuth();
  const isLoggedIn = !!user;
  const userRole = user?.role || "";

  const handleClick = () => {
    if (onEntradaClick) {
      onEntradaClick(id);
    }
  };

  const [showDeleteModal, setShowDeleteModal] = useState(false);

  const handleDelete = () => {
    setShowDeleteModal(true);
  };

  const confirmDelete = () => {
    onDelete(id);
    setShowDeleteModal(false);
  };

  return (
    <Card sx={{ mb: 2 }}>
      <CardActionArea
        component={Link}
        to={`/entrada/${id}`}
        onClick={handleClick}
      >
        <CardContent>
          <Typography variant="h5" component="div">
            {title}
          </Typography>
          <Grid container spacing={2} sx={{ mt: 1 }}>
            <Grid xs={6}>
              <Typography variant="subtitle1" color="textSecondary">
                Creado
              </Typography>
              <Typography variant="body2">
                {new Date(createdAt).toLocaleString("es-ES", {
                  year: "numeric",
                  month: "short",
                  day: "numeric",
                  hour: "2-digit",
                  minute: "2-digit",
                })}
              </Typography>
            </Grid>
          </Grid>
        </CardContent>
      </CardActionArea>
      <Grid xs={6}>
        <Typography variant="subtitle1" color="textSecondary">
          Autor
        </Typography>
        <Typography
          component={Link}
          to={`/perfil/${author.id}`}
          variant="body2"
        >
          Autor: {author.name}
        </Typography>
      </Grid>
      {isLoggedIn && userRole !== "redactor" && (
        <Grid>
          <IconButton color="error" onClick={handleDelete}>
            <DeleteIcon />
          </IconButton>
        </Grid>
      )}
      <ConfirmationModal
        show={showDeleteModal}
        handleClose={() => setShowDeleteModal(false)}
        handleConfirm={confirmDelete}
        message="¿Estás seguro de que deseas eliminar esta entrada?"
      />
    </Card>
  );
};

EntradaCard.propTypes = {
  id: PropTypes.string.isRequired,
  title: PropTypes.string.isRequired,
  author: PropTypes.object.isRequired,
  createdAt: PropTypes.string.isRequired,
  onEntradaClick: PropTypes.func,
  onDelete: PropTypes.func.isRequired,
};

export default EntradaCard;
