import { useEffect, useState, useRef } from "react";
import { postVersion, getVersion } from "../api/VersionApi.js";
import { useNavigate, useParams } from "react-router-dom";
import {
  Container,
  Paper,
  Typography,
  Button,
  TextField,
  Box,
  IconButton,
  Alert,
  Divider,
  Accordion,
  AccordionSummary,
  AccordionDetails,
} from "@mui/material";
import Grid from "@mui/joy/Grid";
import ReactQuill from "react-quill";
import "react-quill/dist/quill.snow.css";
import { Link } from "react-router-dom";
import { Breadcrumbs } from "@mui/material";
import { getEntry, postEntry } from "../api/EntryApi.js";
import { getWiki } from "../api/WikiApi.js";
import { postMedia, getMedia } from "../api/MediaApi";
import DeleteIcon from "@mui/icons-material/Delete";
import VisibilityIcon from "@mui/icons-material/Visibility";
import VisibilityOffIcon from "@mui/icons-material/VisibilityOff";
import ExpandMoreIcon from "@mui/icons-material/ExpandMore";
import ConfirmationModal from "../components/ConfirmationModal.jsx";

function FormVersionPage() {
  const { entryId, versionId, wikiId } = useParams();
  const [entry, setEntry] = useState({});
  const [wiki, setWiki] = useState({});
  const [version, setVersion] = useState({});
  const [versionError, setVersionError] = useState(null);
  const [formErrors, setFormErrors] = useState({
    editor: "",
    content: "",
    author: "",
    title: "",
  });
  const [uploads, setUploads] = useState([]); // Store selected files
  const [existingImages, setExistingImages] = useState([]);
  const [error, setError] = useState(null);
  const formRef = useRef(null);
  const isNewEntry = wikiId;
  const navigate = useNavigate();

  const MAX_FILE_SIZE = 5 * 1024 * 1024; // 5MB
  const ALLOWED_TYPES = ["image/jpeg", "image/png", "image/gif"];

  // State for the confirmation modal
  const [openConfirmDialog, setOpenConfirmDialog] = useState(false);

  // Fetch the version details
  useEffect(() => {
    if (!isNewEntry) {
      if (versionId) {
        getVersion(versionId)
          .then((data) => {
            if (data && Object.keys(data).length > 0) {
              setVersion(data);
              if (data.media_ids) {
                fetchExistingImages(data.media_ids);
              }
            } else {
              setVersionError("No se encontró la versión solicitada.");
            }
          })
          .catch(() =>
            setVersionError("Se produjo un error al obtener la versión.")
          );
      }
    }
  }, [versionId, isNewEntry]);

  // Fetch the entry details
  useEffect(() => {
    if (!isNewEntry) {
    if (entryId) {
      getEntry(entryId)
        .then((data) => {
          if (data && Object.keys(data).length > 0) {
            setEntry(data);
          } else {
            setVersionError("No se encontró la entrada solicitada.");
          }
        })
        .catch(() =>
          setVersionError("Se produjo un error al obtener la entrada.")
        );
    } else {
      setVersionError("No se proporcionó un ID de entrada válido.");
    }
  }
  }, [entryId, isNewEntry]);

  // Fetch the wiki details
  useEffect(() => {
    if (entry && entry.wiki_id || wikiId) {
      const id = wikiId || entry.wiki_id;
      getWiki(id)
        .then((data) => {
          if (data && Object.keys(data).length > 0) {
            setWiki(data);
          } else {
            setVersionError("No se encontró la wiki asociada a esta entrada.");
          }
        })
        .catch(() =>
          setVersionError("Se produjo un error al obtener la wiki asociada.")
        );
    }
  }, [entry, entryId, wikiId]);

  // Fetch existing images
  const fetchExistingImages = async (mediaIds) => {
    try {
      const mediaPromises = mediaIds.map((id) => getMedia(id));
      const mediaResults = await Promise.all(mediaPromises);
      const imagesWithVisibility = mediaResults.map((image) => ({
        ...image,
        isVisible: true,
      }));
      setExistingImages(imagesWithVisibility);
    } catch (error) {
      console.error("Error fetching existing images:", error);
    }
  };

  // Toggle visibility of existing images
  const handleToggleVisibility = (index) => {
    setExistingImages((prevImages) =>
      prevImages.map((image, i) =>
        i === index ? { ...image, isVisible: !image.isVisible } : image
      )
    );
  };

  // Handler for ReactQuill editor change
  const handleEditorChange = (content) => {
    setVersion((prevVersion) => ({
      ...prevVersion,
      content: content,
    }));
    formErrors.content = "";
  };

  // Handler for TextField change
  const handleVersionChange = (event) => {
    const { name, value } = event.target;
    setVersion((prevVersion) => ({
      ...prevVersion,
      [name]: value,
    }));
    setFormErrors((prevErrors) => ({
      ...prevErrors,
      [name]: "",
    }));
  };

  const handleEntryChange = (event) => {
    const { name, value } = event.target;
    setEntry((prevEntry) => ({
      ...prevEntry,
      [name]: value,
    }));
    setFormErrors((prevErrors) => ({
      ...prevErrors,
      [name]: "",
    }));
  };

  // Handler for image upload (store files without uploading)
  const handleImageChange = (event) => {
    const files = Array.from(event.target.files);
    setError(null);

    const validFiles = files.filter((file) => {
      if (!ALLOWED_TYPES.includes(file.type)) {
        setError(`Tipo de archivo no permitido: ${file.name}`);
        return false;
      }
      if (file.size > MAX_FILE_SIZE) {
        setError(`Archivo demasiado grande: ${file.name}`);
        return false;
      }
      return true;
    });

    if (validFiles.length === 0) {
      return;
    }

    setUploads((prevUploads) => [...prevUploads, ...validFiles]);
  };

  // Handler for removing an image
  const handleRemoveImage = (index, isExisting = false) => {
    if (isExisting) {
      setExistingImages((prevImages) =>
        prevImages.filter((_, i) => i !== index)
      );
    } else {
      setUploads((prevUploads) => prevUploads.filter((_, i) => i !== index));
    }
  };

  // Validation function
  const validate = () => {
    let isValid = true;
    const errors = {
      editor: "",
      content: "",
    };
    if (!version.editor && !isNewEntry) {
      errors.editor = "Introduzca un editor";
      isValid = false;
    }
    if (
      !version.content ||
      version.content.replace(/<[^>]+>/g, "").trim().length === 0
    ) {
      errors.content = "Introduzca contenido";
      isValid = false;
    }
    if (!entry.title) {
      errors.title = "Introduzca un título";
      isValid = false;
    }
    if (!entry.author) {
      errors.author = "Introduzca un autor";
      isValid = false;
    }
    setFormErrors(errors);
    return isValid;
  };

  const handleSubmit = (event) => {
    event.preventDefault();
    if (validate()) {
      // Open the confirmation modal
      setOpenConfirmDialog(true);
    }
  };

  const handleConfirmSubmit = async () => {
    // Close the modal
    setOpenConfirmDialog(false);
    // Proceed with form submission logic

    if (!validate()) {
      return;
    }

    let newMediaIds = [];

    const entryData = {
      title: entry.title,
      wiki_id: wikiId,
      author: entry.author,
    };

    // Upload new images
    if (uploads.length > 0) {
      for (const file of uploads) {
        const formData = new FormData();
        formData.append("image", file);

        try {
          //search for the image in the existing images
          // const existingImage = await searchMedia({ publicId: file.name }); //Quizá habría que buscar dentro de la entrada, para que se puedan subir imagenes en otras entradas con el mismo nombre.
          // if (existingImage.length > 0) {
          //   newMediaIds.push(existingImage[0].id);
          //   continue;
          // }
          // else {
          // const response = await postMedia(formData);
          // newMediaIds.push(response.id);
          // }

          const response = await postMedia(formData);
          newMediaIds.push(response.id);
        } catch (error) {
          console.error(`Error uploading file ${file.name}:`, error);
          setError(`Error uploading file ${file.name}`);
          return;
        }
      }
    }

    const versionData = {
      content: version.content,
      editor: version.editor ? version.editor : entry.author,
      entry_id: entryId,
      address: version.address,
      media_ids: [
        ...existingImages
          .filter((image) => image.isVisible)
          .map((image) => image.id),
        ...newMediaIds,
      ],
    };

    try {
      if (isNewEntry) {
        const newEntry = await postEntry(entryData);
        versionData.entry_id = newEntry.id; // Set the new entry ID
      }
      await postVersion(versionData);
      navigate(`/entrada/${versionData.entry_id}`);
    } catch (error) {
      console.error("Error posting version:", error);
    }
  };

  const handleCancelSubmit = () => {
    // Close the modal without submitting
    setOpenConfirmDialog(false);
  };

  return (
    <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
      <Breadcrumbs sx={{ mb: 2 }}>
        <Typography className="breadcrumb-link" component={Link} to="/">
          Inicio
        </Typography>
        <Typography
          className="breadcrumb-link"
          component={Link}
          to={`/wiki/${wiki.id}`}
        >
          {wiki.title}
        </Typography>
        {!isNewEntry && (
        <Typography
          className="breadcrumb-link"
          component={Link}
          to={`/entrada/${entry.id}`}
        >
          {entry.title}
        </Typography>
        )}
      </Breadcrumbs>

      <Paper elevation={3} sx={{ p: { xs: 2, md: 4 }, mb: 4 }}>
        <form id="miFormulario" ref={formRef} onSubmit={handleSubmit}>
          <Grid container spacing={2} alignItems="center">
            <Grid xs={12} sm={8}>
              <Typography variant="h4">
                {isNewEntry ? "Crear Entrada" : "Editar entrada"}
              </Typography>
            </Grid>
            <Grid xs={12} sm={4}>
              <TextField
                id={isNewEntry ? "author" : "editor"}
                name={isNewEntry ? "author" : "editor"}
                label={isNewEntry ? "Autor *" : "Editor *"}
                value={isNewEntry ? entry.author || "" : version.editor || ""}
                onChange={isNewEntry ? handleEntryChange : handleVersionChange}
                variant="outlined"
                error={isNewEntry ? !!formErrors.author : !!formErrors.editor}
                helperText={isNewEntry ? formErrors.author : formErrors.editor}
                fullWidth
                slotProps={{
                  inputLabel: {
                    shrink: true,
                  },
                }}
              />
            </Grid>
            {isNewEntry && (
              <Grid xs={12} sm={4}>
                <TextField
                  id="title"
                  name="title"
                  label="Título *"
                  value={entry.title || ""}
                  onChange={handleEntryChange}
                  variant="outlined"
                  error={!!formErrors.title}
                  helperText={formErrors.title}
                  fullWidth
                  slotProps={{
                    inputLabel: {
                      shrink: true,
                    },
                  }}
                />
              </Grid>
            )}
            <Grid xs={12} sm={4}>
              <TextField
                id="address"
                name="address"
                label="Ubicación (opcional)"
                value={version.address || ""}
                onChange={handleVersionChange}
                variant="outlined"
                fullWidth
                slotProps={{
                  inputLabel: {
                    shrink: true,
                  },
                }}
              />
            </Grid>
          </Grid>
          {versionError &&
            (
              <Typography variant="body1" color="error" gutterBottom>
                {versionError}
              </Typography>
            )}
          <Typography variant="subtitle1" sx={{ mt: 3, mb: 1 }}>
            Contenido:
          </Typography>
          <Box
            sx={{
              height: "400px",
            }}
          >
            <ReactQuill
              theme="snow"
              value={version.content || ""}
              onChange={handleEditorChange}
              style={{
                height: "100%",
                display: "flex",
                flexDirection: "column",
                border: formErrors.content ? "1px solid red" : "1px solid #ccc",
              }}
            />
            {formErrors.content && (
              <Typography
                variant="body2"
                sx={{ color: "red", fontSize: "12px", mt: 1 }}
              >
                {formErrors.content}
              </Typography>
            )}
          </Box>
          <br /> <br />
          <Button
            component="label"
            sx={{ mt: 2 }}
            variant="outlined"
            color="primary"
          >
            Añadir Imágenes
            <input
              id="image-input"
              type="file"
              hidden
              accept="image/*"
              multiple
              onChange={handleImageChange}
            />
          </Button>
          {uploads.length > 0 && (
            <>
              <Divider sx={{ my: 2 }} />
              <Box sx={{ mt: 2 }}>
                <Typography variant="subtitle1">Nuevas Imágenes:</Typography>
                {uploads.map((file, index) => (
                  <Box
                    key={index}
                    sx={{ display: "flex", alignItems: "center", mt: 1 }}
                  >
                    <Typography variant="body2">{file.name}</Typography>
                    <IconButton
                      onClick={() => handleRemoveImage(index)}
                      sx={{ ml: 1 }}
                    >
                      <DeleteIcon />
                    </IconButton>
                  </Box>
                ))}
              </Box>
            </>
          )}
          {existingImages.length > 0 && (
            <Accordion sx={{ mt: 2 }}>
              <AccordionSummary
                expandIcon={<ExpandMoreIcon />}
                aria-controls="panel1a-content"
                id="panel1a-header"
              >
                <Typography>Imágenes Existentes</Typography>
              </AccordionSummary>
              <AccordionDetails>
                {existingImages.map((image, index) => (
                  <Box
                    key={image.id || index}
                    sx={{ display: "flex", alignItems: "center", mt: 1 }}
                  >
                    <a
                      href={image.uploadUrl}
                      target="_blank"
                      rel="noopener noreferrer"
                    >
                      <Typography
                        variant="body2"
                        color="primary"
                        sx={{ textDecoration: "underline", cursor: "pointer" }}
                      >
                        {image.publicId}
                      </Typography>
                    </a>
                    <IconButton
                      onClick={() => handleToggleVisibility(index)}
                      sx={{ ml: 1 }}
                    >
                      {image.isVisible ? (
                        <VisibilityIcon />
                      ) : (
                        <VisibilityOffIcon />
                      )}
                    </IconButton>
                  </Box>
                ))}
              </AccordionDetails>
            </Accordion>
          )}
          {error && <Alert severity="error">{error}</Alert>}
          <Box sx={{ mt: 5, pt: 1 }} display="flex" justifyContent="flex-end">
            <Button type="submit" variant="contained" color="primary">
              Enviar
            </Button>
          </Box>
        </form>
      </Paper>
      {/* Confirmation Modal */}
      <ConfirmationModal
        show={openConfirmDialog}
        handleClose={handleCancelSubmit}
        handleConfirm={handleConfirmSubmit}
        message="¿Está seguro de que desea enviar el formulario?"
      />
    </Container>
  );
}

export default FormVersionPage;
