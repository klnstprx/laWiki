import PropTypes from "prop-types";
import { Link } from "react-router-dom";
import { Card, CardContent, Typography, Grid2, IconButton } from "@mui/material";
import DeleteIcon from "@mui/icons-material/Delete";
import ConfirmationModal from '../components/ConfirmationModal.jsx';
import { useState } from 'react';

const VersionCard = ({ entradaId, versionId, editor, created_at, onDelete }) => {
  const [showDeleteModal, setShowDeleteModal] = useState(false);

  const handleDelete = () => {
    setShowDeleteModal(true);
  }; 

  const confirmDelete = () => {
    onDelete(versionId);
    setShowDeleteModal(false);
  };

  return (
    <Card sx={{ mb: 2 }}>
      <CardContent>
        <Grid2 container spacing={2} alignItems="center">
          <Grid2 xs={12} sm={5}>
            <Typography variant="body1">
              <strong>Fecha:</strong>{" "}
              {new Date(created_at).toLocaleString('es-ES', {
                year: 'numeric',
                month: 'short',
                day: 'numeric',
                hour: '2-digit',
                minute: '2-digit',
              })}
            </Typography>
          </Grid2>
          <Grid2 xs={12} sm={5}>
            <Typography variant="body1">
              <strong>Editor:</strong> {editor}
            </Typography>
          </Grid2>
          <Grid2 xs={12} sm={2}>
            <Typography variant="body1">
              <Link to={`/entrada/${entradaId}/${versionId}`}>Ver</Link>
            </Typography>
          </Grid2>
          <Grid2>
            <IconButton color="error" onClick={handleDelete}>
                <DeleteIcon />
            </IconButton>
          </Grid2>
        </Grid2>
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
