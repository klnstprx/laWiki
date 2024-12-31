import PropTypes from "prop-types";
import { Link } from "react-router-dom";
import { Card, CardContent, IconButton, Typography } from "@mui/material";
import Grid from "@mui/joy/Grid";
import DeleteIcon from "@mui/icons-material/Delete";
import ConfirmationModal from "../components/ConfirmationModal.jsx";
import { useState } from "react";

const VersionCard = ({
  entradaId,
  versionId,
  editor,
  created_at,
  onDelete,
}) => {
  const [showDeleteModal, setShowDeleteModal] = useState(false);
  const isLoggedIn = !!sessionStorage.getItem("appUser"); // Verifica si el usuario está logueado

  const handleDelete = () => {
    setShowDeleteModal(true);
  };

  const confirmDelete = () => {
    onDelete(versionId);
    setShowDeleteModal(false);
  };

  return (
    <Card
      sx={{
        mb: 2,
        "&:hover": {
          boxShadow: 6,
        },
        transition: "box-shadow 0.3s",
      }}
    >
      <CardContent>
        <Grid container spacing={2} alignItems="center">
          <Grid item xs={12} sm={5}>
            <Typography variant="body1">
              <strong>Fecha:</strong>{" "}
              {new Date(created_at).toLocaleString("es-ES", {
                year: "numeric",
                month: "short",
                day: "numeric",
                hour: "2-digit",
                minute: "2-digit",
              })}
            </Typography>
          </Grid>
          <Grid item xs={12} sm={5}>
            <Typography variant="body1">
              <strong>Editor:</strong> {editor}
            </Typography>
          </Grid>
          <Grid item xs={12} sm={2}>
            <Typography variant="body1">
              <Link to={`/entrada/${entradaId}/${versionId}`}>Ver</Link>
            </Typography>
          </Grid>
          {isLoggedIn && sessionStorage.getItem("role") !== "redactor" && (
            <Grid item xs="auto">
              <IconButton color="error" onClick={handleDelete}>
                <DeleteIcon />
              </IconButton>
            </Grid>
          )}
        </Grid>
      </CardContent>
      <ConfirmationModal
        show={showDeleteModal}
        handleClose={() => setShowDeleteModal(false)}
        handleConfirm={confirmDelete}
        message="¿Estás seguro de que deseas eliminar esta versión?"
      />
    </Card>
  );
};

VersionCard.propTypes = {
  entradaId: PropTypes.string.isRequired,
  versionId: PropTypes.string.isRequired,
  editor: PropTypes.string.isRequired,
  created_at: PropTypes.string.isRequired,
  onDelete: PropTypes.func.isRequired,
};

export default VersionCard;
