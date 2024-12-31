import PropTypes from "prop-types";
import { Link } from "react-router-dom";
import { Card, CardContent, IconButton, Typography } from "@mui/material";
import Grid from "@mui/joy/Grid";
import DeleteIcon from "@mui/icons-material/Delete";
import ConfirmationModal from "../components/ConfirmationModal.jsx";
import { useEffect, useState } from "react";
import { getUser } from "../api/AuthApi";

const VersionCard = ({
  entradaId,
  versionId,
  editorId,
  created_at,
  onDelete,
}) => {
  const [showDeleteModal, setShowDeleteModal] = useState(false);
  // Retrieve the JSON string for 'appUser' from sessionStorage
  const appUserJson = sessionStorage.getItem("appUser");

  // Check if 'appUser' exists in sessionStorage
  const isLoggedIn = !!appUserJson;

  let role = null;

  if (isLoggedIn) {
    const appUser = JSON.parse(appUserJson);

    role = appUser.role;
  }

  const handleDelete = () => {
    setShowDeleteModal(true);
  };

  const confirmDelete = () => {
    onDelete(versionId);
    setShowDeleteModal(false);
  };

  const [editor, setEditor] = useState({}); // add state

  useEffect(() => {
    const fetchUsuario = async () => {
      try {
        const userData = await getUser(editorId);
        setEditor(userData);
      } catch (error) {
        console.error("Error fetching user:", error);
      }
    };

    if (editorId) {
      fetchUsuario();
    }
  }, [editorId]);

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
              Editor: <a href={`/perfil/${editor.id}`}>{editor.name}</a>
            </Typography>
          </Grid>
          <Grid item xs={12} sm={2}>
            <Typography variant="body1">
              <Link to={`/entrada/${entradaId}/${versionId}`}>Ver</Link>
            </Typography>
          </Grid>
          {isLoggedIn && role !== "redactor" && (
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
  editorId: PropTypes.string.isRequired,
  created_at: PropTypes.string.isRequired,
  onDelete: PropTypes.func.isRequired,
};

export default VersionCard;
