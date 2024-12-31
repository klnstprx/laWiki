import { Link } from "react-router-dom";
import PropTypes from "prop-types";
import {
  Card,
  CardActionArea,
  CardContent,
  IconButton,
  Typography,
} from "@mui/material";
import DeleteIcon from "@mui/icons-material/Delete";
import Grid from "@mui/joy/Grid";
import ConfirmationModal from "../components/ConfirmationModal.jsx";
import { useEffect, useState } from "react";
import { getUser } from "../api/AuthApi";

const EntradaCard = ({
  id,
  title,
  author,
  createdAt,
  onEntradaClick,
  onDelete,
}) => {
  const handleClick = () => {
    if (onEntradaClick) {
      onEntradaClick(id);
    }
  };

  const [showDeleteModal, setShowDeleteModal] = useState(false);
  const [usuario, setUsuario] = useState({});
  const isLoggedIn = !!sessionStorage.getItem("appUser");

  //cargar usuario de la base de datos
  useEffect(() => {
    const fetchUser = async () => {
      const user = await getUser(author);
      setUsuario(user);
    };
    fetchUser();
  }, [author]);

  const handleDelete = () => {
    setShowDeleteModal(true);
  };

  const confirmDelete = () => {
    onDelete(id);
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
      {/* CardActionArea wraps the content that links to the entry */}
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

      {/* Author info and delete button are outside CardActionArea */}
      <CardContent sx={{ pt: 0 }}>
        <Grid container alignItems="center" justifyContent="space-between">
          <Grid xs={12} sm={6}>
            <Typography variant="subtitle1" color="textSecondary">
              Autor
            </Typography>
            <Typography variant="body2">
              <Link to={`/perfil/${usuario.id}`}>{usuario.name}</Link>
            </Typography>
          </Grid>
          {isLoggedIn && sessionStorage.getItem("role") !== "redactor" && (
            <Grid>
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
  onDelete: PropTypes.func.isRequired,
};

export default EntradaCard;
