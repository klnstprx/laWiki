import { Link } from "react-router-dom";
import PropTypes from "prop-types";
import {
  Card,
  CardActionArea,
  CardContent,
  Typography,
  Grid2,
  IconButton
} from "@mui/material";
import DeleteIcon from "@mui/icons-material/Delete";
import ConfirmationModal from '../components/ConfirmationModal.jsx'; // add import
import { useState } from 'react'; // add import

const EntradaCard = ({ id, title, author, createdAt, onEntradaClick, onDelete }) => {
  const handleClick = () => {
    if (onEntradaClick) {
      onEntradaClick(id);
    }
  };

  const [showDeleteModal, setShowDeleteModal] = useState(false); // add state

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
          <Grid2 container spacing={2} sx={{ mt: 1 }}>
            <Grid2 xs={6}>
              <Typography variant="subtitle1" color="textSecondary">
                Autor
              </Typography>
              <Typography variant="body2">{author}</Typography>
            </Grid2>
            <Grid2 xs={6}>
              <Typography variant="subtitle1" color="textSecondary">
                Creado
              </Typography>
              <Typography variant="body2">{new Date(createdAt).toLocaleString('es-ES', {
                year: 'numeric',
                month: 'short',
                day: 'numeric',
                hour: '2-digit',
                minute: '2-digit',
              })}</Typography>
            </Grid2>
          </Grid2>
        </CardContent>
      </CardActionArea>
      <Grid2>
        <IconButton color="error" onClick={handleDelete}>
            <DeleteIcon />
        </IconButton>
      </Grid2>
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
  author: PropTypes.string.isRequired,
  createdAt: PropTypes.string.isRequired,
  onEntradaClick: PropTypes.func,
  onDelete: PropTypes.func.isRequired
};

export default EntradaCard;
