import { useEffect, useState, useRef, useCallback } from "react";
import {
  searchComments,
  postComment,
  deleteComment,
} from "../api/CommentApi.js";
import { getEntry } from "../api/EntryApi.js";
import { getMedia } from "../api/MediaApi.js";
import { getVersion, searchVersions } from "../api/VersionApi.js";
import { useParams, Link } from "react-router-dom";
import Comentario from "../components/Comentario.jsx";
import Version from "../components/Version.jsx";
import ConfirmationModal from "../components/ConfirmationModal.jsx";
import { useToast } from "../context/ToastContext.jsx";
import {
  Stack,
  Container,
  Paper,
  Typography,
  Button,
  Alert,
  TextField,
  Divider,
  Rating,
  Tooltip,
  IconButton,
} from "@mui/material";
import InfoIcon from "@mui/icons-material/Info";
import HistoryIcon from "@mui/icons-material/History";
import EditIcon from "@mui/icons-material/Edit";
import PersonIcon from "@mui/icons-material/Person";
import CalendarTodayIcon from "@mui/icons-material/CalendarToday";
import { Carousel } from "react-responsive-carousel";
import "react-responsive-carousel/lib/styles/carousel.min.css"; // Import the carousel styles

import Grid from "@mui/joy/Grid";

function EntradaPage() {
  const { entryId, versionId } = useParams();
  const [entry, setEntry] = useState({});
  const [version, setVersion] = useState({});
  const [comments, setComments] = useState([]);
  const [mediaList, setMediaList] = useState([]);
  const [mediaError, setMediaError] = useState(null);
  const [entryError, setEntryError] = useState(null);
  const [commentsError, setCommentsError] = useState(null);
  const [versionError, setVersionError] = useState(null);
  const [coordinates, setCoordinates] = useState(null);

  const [showModal, setShowModal] = useState(false);
  const [pendingComment, setPendingComment] = useState(null);
  const { showToast } = useToast();
  const formRef = useRef(null);

  const [actualVersionId, setActualVersionId] = useState(versionId || null);

  const geoCacheRef = useRef(
    JSON.parse(sessionStorage.getItem("geoCache")) || {},
  ); // cache de geocoding

  const saveCacheToSessionStorage = () => {
    sessionStorage.setItem("geoCache", JSON.stringify(geoCacheRef.current));
  };

  const fetchCoordinatesNominatim = useCallback(async (address) => {
    // Comprueba si la dirección ya está en el cache
    if (geoCacheRef.current[address]) {
      console.log(
        "Obteniendo coordenadas desde el cache en memoria:",
        geoCacheRef.current[address],
      );
      return geoCacheRef.current[address]; // Retorna las coordenadas almacenadas
    }

    // Si no está en el cache, realiza la solicitud a la API
    const url = `https://nominatim.openstreetmap.org/search?q=${encodeURIComponent(
      address,
    )}&format=json&addressdetails=1&limit=1`;

    try {
      const response = await fetch(url);
      const data = await response.json();

      if (data.length > 0) {
        const { lat, lon } = data[0];
        const coordinates = { lat: parseFloat(lat), lon: parseFloat(lon) };
        console.log("Coordenadas obtenidas de la API:", coordinates);

        // Almacena las coordenadas en el cache y en sessionStorage antes de retornarlas
        geoCacheRef.current[address] = coordinates;
        saveCacheToSessionStorage(); // Actualiza sessionStorage
        return coordinates;
      } else {
        throw new Error("No se encontraron coordenadas para la dirección.");
      }
    } catch (error) {
      console.error("Error al realizar la geocodificación:", error);
      return null;
    }
  }, []);

  // Handler to close the confirmation modal
  const handleClose = () => {
    setShowModal(false);
    setPendingComment(null);
    showToast("El comentario no se ha creado", "warning");
  };

  // Handler to confirm and post the comment
  const handleConfirm = async () => {
    setShowModal(false);
    try {
      const result = await postComment(pendingComment);
      setComments((prevComments) => [...prevComments, result]);
      formRef.current.reset();
      setPendingComment(null);
      showToast("El comentario se ha creado correctamente!", "success");
    } catch (error) {
      console.error("Error al enviar:", error);
      showToast("Error al enviar el comentario", "error");
    }
  };

  const [showDeleteCommentModal, setShowDeleteCommentModal] = useState(false); // add state
  const [commentToDelete, setCommentToDelete] = useState(null); // add state

  const handleDeleteComment = (commentId) => {
    setCommentToDelete(commentId);
    setShowDeleteCommentModal(true);
  };

  const confirmDeleteComment = async () => {
    await deleteComment(commentToDelete);
    setComments((prevComments) =>
      prevComments.filter((comment) => comment.id !== commentToDelete),
    );
    setShowDeleteCommentModal(false);
    showToast("Comentario eliminado correctamente", "success");
  };

  // Fetch the entry details
  useEffect(() => {
    if (entryId) {
      getEntry(entryId)
        .then((data) => {
          if (data && Object.keys(data).length > 0) {
            setEntry(data);
            fetchMedia(data.media_ids);
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
  }, [entryId]);

  //fetch media
  const fetchMedia = async (mediaIdsArray) => {
    if (!Array.isArray(mediaIdsArray) || mediaIdsArray.length === 0) {
      console.log("No media IDs found or mediaIdsArray is not an array.");
      return;
    }
    try {
      const mediaPromises = mediaIdsArray.map((id) => getMedia(id));
      const mediaResults = await Promise.all(mediaPromises);
      setMediaList(mediaResults);
    } catch (error) {
      console.error("Error fetching media:", error);
      setMediaError("Failed to fetch media");
    }
  };

  // Fetch the version details
  useEffect(() => {
    if (versionId) {
      // If versionId is provided in the URL, fetch that specific version
      setActualVersionId(versionId);
    } else if (entryId) {
      // If no versionId is provided, fetch the latest version for the entry
      searchVersions({ entryID: entryId })
        .then((versions) => {
          if (versions && versions.length > 0) {
            // Sort the versions by createdAt descending to get the latest
            versions.sort(
              (a, b) => new Date(b.created_at) - new Date(a.created_at),
            );
            const latestVersion = versions[0];
            setActualVersionId(latestVersion.id);
          } else {
            setLoadingVersion(false);
            setVersionError(
              "No se encontró ninguna versión para esta entrada.",
            );
          }
        })
        .catch(() =>
          setVersionError("Se produjo un error al obtener las versiones."),
        );
    }
  }, [entryId, versionId]);

  const [loadingVersion, setLoadingVersion] = useState(true);

  // Fetch the version data when actualVersionId changes
  useEffect(() => {
    if (actualVersionId) {
      getVersion(actualVersionId)
        .then(async (data) => {
          if (data && Object.keys(data).length > 0) {
            setVersion(data);
            setLoadingVersion(false); // Data is now loaded
            if (data.address) {
              const coords = await fetchCoordinatesNominatim(data.address);
              setCoordinates(coords);
            }
          } else {
            setVersionError("No se encontró la versión solicitada.");
            setLoadingVersion(false);
          }
        })
        .catch(() => {
          setVersionError("Se produjo un error al obtener la versión.");
          setLoadingVersion(false);
        });

      // Fetch comments for the actual version
      searchComments({ versionID: actualVersionId })
        .then((data) => {
          if (data && data.length > 0) {
            setComments(data);
          } else {
            setComments([]);
          }
        })
        .catch(() =>
          setCommentsError("Se produjo un error al obtener los comentarios."),
        );
    }
  }, [actualVersionId, fetchCoordinatesNominatim]);

  // Handler to submit a new comment
  async function subirComentario(event) {
    event.preventDefault();
    const formData = new FormData(event.target);
    const jsonData = Object.fromEntries(formData.entries());
    jsonData["version_id"] = actualVersionId;
    jsonData["entry_id"] = entryId;
    jsonData["rating"] = parseInt(jsonData["rating"], 10);
    setPendingComment(jsonData);
    setShowModal(true);
  }

  return (
    <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
      {/* Entry Title */}
      {!entryError && entry && (
        <Typography variant="h3" gutterBottom>
          {entry.title}
        </Typography>
      )}

            {/* Entry Details */}
      <Paper elevation={3} sx={{ p: 2, mb: 4 }}>
        {!entryError && entry && (
          <>
            <Tooltip
              title={
                <Typography variant="subtitle1">
                  Autor: {entry.author}
                </Typography>
              }
              arrow
            >
              <IconButton>
                <PersonIcon />
              </IconButton>
            </Tooltip>
            <Tooltip
              title={
                <Typography variant="subtitle1">
                  Fecha de creación:{" "}
                  {new Date(entry.created_at).toLocaleString("es-ES", {
                    year: "numeric",
                    month: "short",
                    day: "numeric",
                    hour: "2-digit",
                    minute: "2-digit",
                  })}
                </Typography>
              }
              arrow
            >
              <IconButton>
                <CalendarTodayIcon />
              </IconButton>
            </Tooltip>
            <Tooltip 
            title={
            <Typography variant="subtitle1">
              Ver historial
            </Typography>}
             arrow>
              <IconButton component={Link} to={`/versiones/${entry.id}/`}>
                <HistoryIcon />
              </IconButton>
            </Tooltip>
            <Tooltip 
            title={
              <Typography variant="subtitle1">
                Editar contenido
              </Typography>
            }
             arrow>
              <IconButton component={Link} to={`/version/form/${entry.id}/${actualVersionId || ""}`}>
                <EditIcon />
              </IconButton>
            </Tooltip>
          </>
        )}
      </Paper>

      {/* Version Content */}
      <Paper elevation={3} sx={{ p: 2, mb: 4 }}>
        {loadingVersion ? (
          <Typography variant="body1">Cargando versión...</Typography>
        ) : versionError ? (
          <Alert severity="error">{versionError}</Alert>
        ) : (
          <Version
            content={version.content}
            editor={version.editor}
            created_at={version.created_at}
            entry_id={version.entry_id}
            address={version.address}
            coordinates={coordinates}
          />
        )}
      </Paper>

      {/* Divider */}
      <Divider sx={{ my: 4 }} />

      {/* Media */}
      {entryError && <Alert severity="error">{entryError}</Alert>}
      {mediaError && <Alert severity="error">{mediaError}</Alert>}
      {mediaList.length > 0 && ( // Verifica si hay elementos en mediaList
        <Paper elevation={3} sx={{ p: 2, mb: 4 }}>
          <Carousel showThumbs={false} infiniteLoop useKeyboardArrows>
            {mediaList.map((media, index) => (
              <div key={index}>
                <img
                  src={media.uploadUrl}
                  alt={media.publicId}
                  style={{ maxWidth: "33%" }}
                />
              </div>
            ))}
          </Carousel>
        </Paper>
      )}

      {/* Entry Details */}
      <Paper elevation={3} sx={{ p: 2, mb: 4 }}>
        {entryError && <Alert severity="error">{entryError}</Alert>}
        {!entryError && entry && (
          <>
            <Typography variant="subtitle1" gutterBottom>
              Autor: {entry.author} | Fecha de creación:{" "}
              {new Date(entry.created_at).toLocaleString("es-ES", {
                year: "numeric",
                month: "short",
                day: "numeric",
                hour: "2-digit",
                minute: "2-digit",
              })}
            </Typography>
            <Typography variant="subtitle1" gutterBottom>
              <Link to={`/versiones/${entry.id}/`}>Ver historial</Link>
              {" | "}
              <Link to={`/version/form/${entry.id}/${actualVersionId || ""}`}>
                Editar contenido
              </Link>
            </Typography>
          </>
        )}
      </Paper>

      {/* Comments */}
      <Paper elevation={3} sx={{ p: 2, mb: 4 }}>
        <Typography variant="h5" gutterBottom>
          Comentarios
        </Typography>
        {commentsError && <Alert severity="error">{commentsError}</Alert>}
        {!commentsError && comments.length > 0 ? (
          <Stack spacing={2} sx={{ mb: 2 }}>
            {comments.map((comment) => (
              <Comentario
                key={comment.id}
                id={comment.id}
                content={comment.content}
                rating={comment.rating}
                created_at={comment.created_at}
                author={comment.author}
                onDelete={(id) => handleDeleteComment(id)}
              />
            ))}
          </Stack>
        ) : (
          !commentsError && (
            <Alert severity="info">No se encontraron comentarios.</Alert>
          )
        )}
      </Paper>

      {/* Form to Add Comment */}
      <Paper elevation={3} sx={{ p: 2, mb: 4 }}>
        <Typography variant="h5" gutterBottom>
          Añadir comentario
        </Typography>
        <form id="miFormulario" ref={formRef} onSubmit={subirComentario}>
          <Grid container spacing={2}>
            <Grid xs={12}>
              <TextField
                id="content"
                name="content"
                label="Contenido"
                multiline
                required
                fullWidth
                rows={4}
              />
            </Grid>
            <Grid xs={12} sm={6} md={4}>
              <Typography variant="subtitle1" gutterBottom>
                Valoración:
              </Typography>
              <Rating name="rating" id="rating" size="large" />
            </Grid>
            <Grid xs={12} sm={6} md={4}>
              <TextField
                id="author"
                name="author"
                label="Autor"
                required
                fullWidth
              />
            </Grid>
            <Grid xs={12} md={4}>
              <Button
                type="submit"
                variant="contained"
                color="primary"
                fullWidth
                sx={{ height: "100%" }}
              >
                Enviar
              </Button>
            </Grid>
          </Grid>
        </form>
      </Paper>

      <ConfirmationModal
        message="¿Estás seguro de que quieres crear este comentario?"
        show={showModal}
        handleClose={handleClose}
        handleConfirm={handleConfirm}
      />

      <ConfirmationModal
        show={showDeleteCommentModal}
        handleClose={() => setShowDeleteCommentModal(false)}
        handleConfirm={confirmDeleteComment}
        message="¿Estás seguro de que deseas eliminar este comentario?"
      />
    </Container>
  );
}

export default EntradaPage;
