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
  const [entrada, setEntrada] = useState(null);
  const [version, setVersion] = useState(null);
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
      // Actualizar el estado para eliminar el comentario borrado
      setComentarios((prevComentarios) =>
        prevComentarios.filter((comment) => comment.id !== commentId),
      );
      showToast("Comentario eliminado correctamente", "success");
    } catch (error) {
      console.error("Error al eliminar el comentario:", error);
      showToast("Error al eliminar el comentario", "danger");
    }
  };

  // Obtener la entrada
  useEffect(() => {
    if (id) {
      getEntry(id)
        .then((data) => {
          if (data && Object.keys(data).length > 0) {
            setEntrada(data);
          } else {
            setEntryError("No se encontró la entrada solicitada.");
          }
        })
        .catch(() =>
          setEntryError("Se produjo un error al obtener la entrada."),
        );
    } else {
      setEntryError("No se proporcionó un ID de entrada válido.");
    }
  }, [id]);

  // Obtener la versión
  useEffect(() => {
    if (versionID) {
      getVersion(versionID)
        .then((data) => {
          if (data && Object.keys(data).length > 0) {
            setVersion(data);
          } else {
            setVersionError("No se encontró la versión solicitada.");
          }
        })
        .catch(() =>
          setVersionError("Se produjo un error al obtener la versión."),
        );
    } else {
      setVersionError("No se proporcionó un ID de versión válido.");
    }
  }, [versionID]);

  // Obtener los comentarios
  useEffect(() => {
    if (versionID) {
      searchComments({ versionID: versionID })
        .then((data) => {
          if (data) {
            setComentarios(data);
          } else {
            setCommentsError(
              "No se encontraron comentarios para esta versión.",
            );
          }
        })
        .catch(() =>
          setCommentsError("Se produjo un error al obtener los comentarios."),
        );
    } else {
      setCommentsError("No se proporcionó un ID de versión válido.");
    }
  }, [versionID]);

  // Handler para enviar el comentario
  async function subirComentario(event) {
    event.preventDefault();

    const formData = new FormData(event.target);
    const jsonData = Object.fromEntries(formData.entries());

    jsonData["version_id"] = versionID;
    jsonData["rating"] = parseInt(jsonData["rating"], 10);

    setPendingComment(jsonData);
    setShowModal(true);
  }

  return (
    <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
      {/* Detalles de la Entrada */}
      <Paper elevation={3} sx={{ p: 2, mb: 4 }}>
        <Typography variant="h4" gutterBottom>
          Datos de la Entrada
        </Typography>
        {entryError && <Alert severity="error">{entryError}</Alert>}
        {!entryError && entrada && (
          <>
            <Typography variant="h6">Título: {entrada.title}</Typography>
            <Typography variant="h6">Autor: {entrada.author}</Typography>
            <Typography variant="h6">
              Fecha de creación:{" "}
              {new Date(entrada.created_at).toLocaleDateString()}
            </Typography>
            <Typography variant="h6">
              <a
                href={`http://localhost:5173/versiones?entry_id=${entrada.id}`}
              >
                Ver historial
              </a>
            </Typography>
          </>
        )}
      </Paper>

      {/* Contenido de la Versión */}
      <Paper elevation={3} sx={{ p: 2, mb: 4 }}>
        <Typography variant="h5" gutterBottom>
          Contenido de la Versión
        </Typography>
        {versionError && <Alert severity="error">{versionError}</Alert>}
        {!versionError && version && (
          <Version
            content={version.content}
            editor={version.editor}
            created_at={version.created_at}
            entry_id={version.entry_id}
          />
        )}
      </Paper>

      {/* Comentarios */}
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
          !commentsError && (
            <Alert severity="info">No se encontraron comentarios.</Alert>
          )
        )}
      </Paper>

      {/* Formulario para añadir comentario */}
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
                label="Calificación"
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
      <ConfirmationModal
        message="¿Estás seguro de que quieres crear este comentario?"
        show={showModal}
        handleClose={handleClose}
        handleConfirm={handleConfirm}
      />
    </Container>
  );
}

export default EntradaPage;
