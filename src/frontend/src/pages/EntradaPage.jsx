import { useCallback, useEffect, useRef, useState } from "react";
import {
  deleteComment,
  postComment,
  searchComments,
} from "../api/CommentApi.js";
import { getEntry, translateEntry } from "../api/EntryApi.js";
import { getVersion, searchVersions } from "../api/VersionApi.js";
import { Link, useParams } from "react-router-dom";
import Comentario from "../components/Comentario.jsx";
import Version from "../components/Version.jsx";
import ConfirmationModal from "../components/ConfirmationModal.jsx";
import { useToast } from "../context/ToastContext.jsx";
import {
  Alert,
  Breadcrumbs,
  Button,
  Container,
  Paper,
  Rating,
  Stack,
  TextField,
  Typography,
  Menu,
  MenuItem,
  Divider,
} from "@mui/material";
import { getWiki } from "../api/WikiApi.js";
import HistoryIcon from "@mui/icons-material/History";
import EditIcon from "@mui/icons-material/Edit";
import { getUser } from "../api/AuthApi.js";
import { useLanguage } from "../context/LanguageContext.jsx";

import Grid from "@mui/joy/Grid";
import { availableLanguages } from "../constants/languages.js";

function EntradaPage() {
  const { entryId, versionId } = useParams();
  const [entry, setEntry] = useState(null);
  const [wiki, setWiki] = useState(null);
  const [version, setVersion] = useState(null);
  const [comments, setComments] = useState([]);
  const [entryError, setEntryError] = useState(null);
  const [commentsError, setCommentsError] = useState(null);
  const [versionError, setVersionError] = useState(null);
  const [coordinates, setCoordinates] = useState(null);
  const [usuario, setUsuario] = useState({}); // add state

  const [showModal, setShowModal] = useState(false);
  const [pendingComment, setPendingComment] = useState(null);
  const { showToast } = useToast();
  const formRef = useRef(null);

  const [actualVersionId, setActualVersionId] = useState(versionId || null);
  const isLoggedIn = !!sessionStorage.getItem("user"); // Verifica si el usuario está logueado

  const geoCacheRef = useRef(
    JSON.parse(sessionStorage.getItem("geoCache")) || {}
  ); // cache de geocoding

  const saveCacheToSessionStorage = () => {
    sessionStorage.setItem("geoCache", JSON.stringify(geoCacheRef.current));
  };

  const fetchCoordinatesNominatim = useCallback(async (address) => {
    // Comprueba si la dirección ya está en el cache
    if (geoCacheRef.current[address]) {
      console.log(
        "Obteniendo coordenadas desde el cache en memoria:",
        geoCacheRef.current[address]
      );
      return geoCacheRef.current[address]; // Retorna las coordenadas almacenadas
    }

    // Si no está en el cache, realiza la solicitud a la API
    const url = `https://nominatim.openstreetmap.org/search?q=${encodeURIComponent(
      address
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

  const [showDeleteCommentModal, setShowDeleteCommentModal] = useState(false);
  const [commentToDelete, setCommentToDelete] = useState(null);

  const handleDeleteComment = (commentId) => {
    setCommentToDelete(commentId);
    setShowDeleteCommentModal(true);
  };

  const confirmDeleteComment = async () => {
    await deleteComment(commentToDelete);
    setComments((prevComments) =>
      prevComments.filter((comment) => comment.id !== commentToDelete)
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

            //cargar usuario de la base de datos
            getUser(data.author)
              .then((user) => {
                setUsuario(user);
              })
              .catch(() => setEntryError("No se pudo cargar el autor."));
          } else {
            setEntryError("No se encontró la entrada solicitada.");
          }
        })
        .catch(() =>
          setEntryError("Se produjo un error al obtener la entrada.")
        );
    } else {
      setEntryError("No se proporcionó un ID de entrada válido.");
    }
  }, [entryId]);

  // Fetch the wiki details
  useEffect(() => {
    if (entry?.wiki_id) {
      getWiki(entry.wiki_id)
        .then((data) => {
          if (data && Object.keys(data).length > 0) {
            setWiki(data);
          } else {
            setEntryError("No se encontró la wiki asociada a esta entrada.");
          }
        })
        .catch(() =>
          setEntryError("Se produjo un error al obtener la wiki asociada.")
        );
    }
  }, [entry?.wiki_id]);

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
              (a, b) => new Date(b.created_at) - new Date(a.created_at)
            );
            const latestVersion = versions[0];
            setActualVersionId(latestVersion.id);
          } else {
            setLoadingVersion(false);
          }
        })
        .catch(() =>
          setVersionError("Se produjo un error al obtener las versiones.")
        );
    }
  }, [entryId, versionId]);

  const [loadingVersion, setLoadingVersion] = useState(true);

  // Fetch the version data when actualVersionId changes
  useEffect(() => {
    if (actualVersionId) {
      setLoadingVersion(true);
      getVersion(actualVersionId)
        .then(async (data) => {
          if (data && Object.keys(data).length > 0) {
            setVersion(data);
            if (data.address) {
              const coords = await fetchCoordinatesNominatim(data.address);
              setCoordinates(coords);
            }
          } else {
            setLoadingVersion(false);
          }
        })
        .catch(() => {
          setVersionError("Se produjo un error al obtener la versión.");
          setLoadingVersion(false);
        })
        .finally(() => {
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
          setCommentsError("Se produjo un error al obtener los comentarios.")
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
    jsonData["author"] = sessionStorage.getItem("id");
    setPendingComment(jsonData);
    setShowModal(true);
  }

  // Language selector state and handlers
  const [isModalOpen, setIsModalOpen] = useState(false);
  const { selectedOption, setSelectedOption } = useLanguage();
  const [anchorEl, setAnchorEl] = useState(null);
  const [pendingLanguage, setPendingLanguage] = useState(null);

  const handleDropdownClose = () => {
    setAnchorEl(null);
  };

  const handleOptionSelect = (option) => {
    setPendingLanguage(option);
    handleDropdownClose();
    setIsModalOpen(true);
  };

  const handleTranslateEntry = async () => {
    if (entry.sourceLang !== pendingLanguage) {
      try {
        await translateEntry(entryId, pendingLanguage);
        showToast(`Entrada traducida a ${pendingLanguage} correctamente`, "success");
        
        // Fetch both updated entry and version data
        const updatedEntry = await getEntry(entryId);
        const updatedVersion = await getVersion(actualVersionId);
        
        setEntry(updatedEntry);
        setVersion(updatedVersion);
        setSelectedOption(pendingLanguage);
      } catch (error) {
        console.error("Error al traducir la entrada:", error);
        showToast("Error al traducir la entrada", "error");
      }
      setIsModalOpen(false);
    } else {
      showToast(`Entrada traducida a ${pendingLanguage} correctamente`, "success");
      setSelectedOption(pendingLanguage);
      setIsModalOpen(false);
    }
};

  const getTranslatedField = (field) => {
    if (entry.sourceLang === selectedOption) {
      return entry[field];
    } else {
      return entry.translatedFields?.[selectedOption]?.[field] || entry[field];
    }
  };

  const getTranslatedFieldWiki = (field) => {
    if (entry.sourceLang === selectedOption) {
      return wiki[field];
    } else {
      return wiki.translatedFields?.[selectedOption]?.[field] || wiki[field];
    }
  };

  const getTranslatedFieldVersion = (field) => {
    if (entry.sourceLang === selectedOption) {
      return version[field];
    } else {
      return (
        version.translatedFields?.[selectedOption]?.[field] || version[field]
      );
    }
  };

  return (
    <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
      <Breadcrumbs sx={{ mb: 2 }}>
        <Typography className="breadcrumb-link" component={Link} to="/">
          Inicio
        </Typography>

        {wiki ? (
          <Typography
            className="breadcrumb-link"
            component={Link}
            to={`/wiki/${wiki.id}`}
          >
            {getTranslatedFieldWiki("title")}
          </Typography>
        ) : (
          <Typography className="breadcrumb-link">Cargando wiki...</Typography>
        )}

        {entry ? (
          <Typography className="breadcrumb-active">
            {getTranslatedField("title")}
          </Typography>
        ) : (
          <Typography className="breadcrumb-active">
            Cargando entrada...
          </Typography>
        )}
      </Breadcrumbs>

      {/* Entry Title */}
      {!entryError && entry && (
        <Paper
          elevation={3}
          sx={{ p: 2, mb: 4, textAlign: "center", borderRadius: 1 }}
        >
          <Stack
            direction="row"
            justifyContent="space-between"
            alignItems="center"
          >
            <Typography variant="subtitle2">
              Autor: <a href={`/perfil/${usuario.id}`}>{usuario.name}</a>
            </Typography>
            <Typography variant="caption" color="text.secondary">
              {new Date(entry.created_at).toLocaleDateString()}
            </Typography>
          </Stack>

          <Divider sx={{ my: 2 }} />

          <Typography variant="h2" component="h1" style={{ padding: "16px" }}>
            {getTranslatedField("title")}
          </Typography>
          <Divider sx={{ my: 2 }} />
          {isLoggedIn && (
            <Stack
              direction="row"
              spacing={2}
              justifyContent="center"
              alignItems="center"
            >
              <Button
                variant="contained"
                component={Link}
                to={`/versiones/${entry.id}/`}
                startIcon={<HistoryIcon />}
              >
                Ver historial
              </Button>
              <Button
                variant="contained"
                component={Link}
                to={`/version/form/${entry.id}/${actualVersionId || ""}`}
                startIcon={<EditIcon />}
              >
                Editar contenido
              </Button>
              <Button
                variant="contained"
                color="primary"
                sx={{ mt: 2, mr: 2 }}
                onClick={(e) => {
                  setAnchorEl(e.currentTarget);
                  setPendingLanguage(null);
                }}
              >
                Cambiar Idioma: {selectedOption || "Seleccionar"}
              </Button>
              {["admin", "editor", "redactor"].includes(
                sessionStorage.getItem("role")
              ) && (
                <Button
                  variant="contained"
                  color="secondary"
                  sx={{ mt: 2 }}
                  onClick={(e) => {
                    setAnchorEl(e.currentTarget);
                    setPendingLanguage("translate");
                  }}
                >
                  Traducir
                </Button>
              )}
              <Menu
                anchorEl={anchorEl}
                open={Boolean(anchorEl)}
                onClose={handleDropdownClose}
              >
                {availableLanguages.map((lang) => {
                  if (!pendingLanguage) {
                    if (
                      entry.translatedFields?.[lang.code] ||
                      entry.sourceLang === lang.code
                    ) {
                      return (
                        <MenuItem
                          key={lang.code}
                          onClick={() => {
                            setSelectedOption(lang.code);
                            handleDropdownClose();
                          }}
                        >
                          {lang.name}
                        </MenuItem>
                      );
                    }
                    return null;
                  } else if (entry.sourceLang !== lang.code) {
                    return (
                      <MenuItem
                        key={lang.code}
                        onClick={() => handleOptionSelect(lang.code)}
                      >
                        {lang.name}{" "}
                        {entry.translatedFields?.[lang.code]
                          ? "(Actualizar)"
                          : ""}
                      </MenuItem>
                    );
                  }
                  return null;
                })}
              </Menu>
            </Stack>
          )}
        </Paper>
      )}

      {/* Version Content */}
      <Paper elevation={3} sx={{ p: 2, mb: 4 }}>
        {loadingVersion ? (
          <Typography variant="body1">Cargando versión...</Typography>
        ) : versionError ? (
          <Alert severity="error">{versionError}</Alert>
        ) : !version ? (
          <Alert severity="info">
            No se ha encontrado niguna version asignada a esta entrada.
          </Alert>
        ) : (
          <Version
            content={getTranslatedFieldVersion("content")}
            editor={version.editor}
            created_at={version.created_at}
            entry_id={version.entry_id}
            address={version.address}
            coordinates={coordinates}
            media_ids={version.media_ids}
          />
        )}
      </Paper>

      {entryError && <Alert severity="error">{entryError}</Alert>}

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
      {isLoggedIn && (
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
                  slotProps={{
                    inputLabel: {
                      shrink: true,
                    },
                  }}
                />
              </Grid>
              <Grid xs={12} sm={6} md={4}>
                <Typography variant="subtitle2" gutterBottom>
                  Valoración:
                </Typography>
                <Rating name="rating" id="rating" size="large" />
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
      )}

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

      <ConfirmationModal
        show={isModalOpen}
        handleClose={() => setIsModalOpen(false)}
        handleConfirm={handleTranslateEntry}
        message={`¿Estás seguro de que quieres traducir esta entrada a ${pendingLanguage}?`}
      />
    </Container>
  );
}

export default EntradaPage;
