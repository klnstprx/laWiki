import { useEffect, useState, useRef } from "react";
import {
  searchComments,
  postComment,
  deleteComment,
} from "../api/CommentApi.js";
import { getEntry } from "../api/EntryApi.js";
import { getVersion } from "../api/VersionApi.js";
import { useSearchParams } from "react-router-dom";
import Comentario from "../components/Comentario.jsx";
import Version from "../components/Version.jsx";
import MainLayout from "../layout/MainLayout.jsx";
import ConfirmationModal from "../components/ConfirmationModal.jsx";
import { useToast } from "../context/ToastContext.1.jsx";
import {
  Container,
  Paper,
  Typography,
  Button,
  Alert,
  List,
  ListItem,
  Input,
  Grid2,
} from "@mui/material";

function EntradaPage() {
  const [entrada, setEntrada] = useState({});
  const [version, setVersion] = useState({});
  const [comentarios, setComentarios] = useState([]);
  const [entryError, setEntryError] = useState(null);
  const [commentsError, setCommentsError] = useState(null);
  const [versionError, setVersionError] = useState(null);

  const [showModal, setShowModal] = useState(false);
  const [pendingComment, setPendingComment] = useState(null);
  const { showToast } = useToast();
  const formRef = useRef(null);

  const [searchParams] = useSearchParams();
  const id = searchParams.get("id");
  const versionID = searchParams.get("versionID");

  const handleClose = () => {
    setShowModal(false);
    setPendingComment(null);
    showToast("El comentario no se ha creado", "danger");
  };

  const handleConfirm = async () => {
    setShowModal(false);
    try {
      const result = await postComment(pendingComment);

      setComentarios((prevComentarios) => [...prevComentarios, result]);

      formRef.current.reset();
      setPendingComment(null);

      showToast("El comentario se ha creado correctamente!", "success");
    } catch (error) {
      console.error("Error al enviar:", error);
      showToast("Error al enviar el comentario", "danger");
    }
  };

  const handleDeleteComment = async (commentId) => {
    try {
      await deleteComment(commentId);
      // Update state to remove the deleted comment
      setComentarios((prevComentarios) =>
        prevComentarios.filter((comment) => comment.id !== commentId),
      );
      showToast("Comentario eliminado correctamente", "success");
    } catch (error) {
      console.error("Error al eliminar el comentario:", error);
      showToast("Error al eliminar el comentario", "danger");
    }
  };

  // Fetch data on mount
  useEffect(() => {
    getEntry(id)
      .then(setEntrada)
      .catch((err) => setEntryError(err.message));
  }, [id]);

  useEffect(() => {
    searchComments({ versionID: versionID })
      .then(setComentarios)
      .catch((err) => setCommentsError(err.message));
  }, [versionID]);

  useEffect(() => {
    getVersion(versionID)
      .then(setVersion)
      .catch((err) => setVersionError(err.message));
  }, [versionID]);

  // Form submission handler
  async function subirComentario(event) {
    event.preventDefault();

    const formData = new FormData(event.target);
    const jsonData = Object.fromEntries(formData.entries());

    jsonData["version_id"] = version.id;
    jsonData["rating"] = parseInt(jsonData["rating"], 10);

    // Store the comment data and show the modal
    setPendingComment(jsonData);
    setShowModal(true);
  }

  return (
    <MainLayout>
      <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
        {/* Entry Details */}
        <Paper elevation={3} sx={{ p: 2, mb: 4 }}>
          <Typography variant="h4" gutterBottom>
            Datos de la Wiki
          </Typography>
          {entryError && <Alert severity="error">{entryError}</Alert>}
          {!entryError && (
            <>
              <Typography variant="h6">Título: {entrada.title}</Typography>
              <Typography variant="h6">Autor: {entrada.author}</Typography>
              <Typography variant="h6">
                Fecha de creación:{" "}
                {new Date(entrada.created_at).toLocaleDateString()}
              </Typography>
              <Typography variant="h6">
                <a
                  href={`http://localhost:5173/versiones?entryID=${entrada.id}`}
                >
                  Ver historial
                </a>
              </Typography>
            </>
          )}
        </Paper>

        {/* Version Content */}
        <Paper elevation={3} sx={{ p: 2, mb: 4 }}>
          <Typography variant="h5" gutterBottom>
            Contenido de la versión
          </Typography>
          {versionError && <Alert severity="error">{versionError}</Alert>}
          {!versionError && (
            <Version
              content={version.content}
              editor={version.editor}
              created_at={version.created_at}
              entry_id={version.entry_id}
            />
          )}
        </Paper>

        {/* Comments */}
        <Paper elevation={3} sx={{ p: 2, mb: 4 }}>
          <Typography variant="h5" gutterBottom>
            Comentarios
          </Typography>
          {commentsError && <Alert severity="error">{commentsError}</Alert>}
          {!commentsError && comentarios.length > 0 ? (
            <List>
              {comentarios.map((comentario) => (
                <ListItem key={comentario.id}>
                  <Comentario
                    id={comentario.id}
                    content={comentario.content}
                    rating={comentario.rating}
                    created_at={comentario.created_at}
                    author={comentario.author}
                    onDelete={handleDeleteComment}
                  />
                </ListItem>
              ))}
            </List>
          ) : (
            <Alert severity="info">No se encontraron comentarios.</Alert>
          )}
        </Paper>

        {/* Add Comment Form */}
        <Paper elevation={3} sx={{ p: 2, mb: 4 }}>
          <Typography variant="h5" gutterBottom>
            Añadir comentario
          </Typography>
          <form id="miFormulario" ref={formRef} onSubmit={subirComentario}>
            <Grid2 container spacing={2}>
              <Grid2 item xs={12}>
                <Input
                  id="content"
                  name="content"
                  label="Contenido"
                  multiline
                  required
                  fullWidth
                />
              </Grid2>
              <Grid2 item xs={12} sm={6} md={4}>
                <Input
                  id="rating"
                  name="rating"
                  label="Rating"
                  type="number"
                  inputProps={{ min: 1, max: 5 }}
                  required
                  fullWidth
                />
              </Grid2>
              <Grid2 item xs={12} sm={6} md={4}>
                <Input
                  id="author"
                  name="author"
                  label="Autor"
                  required
                  fullWidth
                />
              </Grid2>
              <Grid2 item xs={12} md={4}>
                <Button
                  type="submit"
                  variant="contained"
                  color="primary"
                  fullWidth
                  sx={{ height: "100%" }}
                >
                  Enviar
                </Button>
              </Grid2>
            </Grid2>
          </form>
        </Paper>
      </Container>

      {/* Confirmation Modal */}
      <ConfirmationModal
        message="¿Estás seguro de que quieres crear este comentario?"
        show={showModal}
        handleClose={handleClose}
        handleConfirm={handleConfirm}
      />
    </MainLayout>
  );
}

export default EntradaPage;
