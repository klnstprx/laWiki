import { useEffect, useRef, useState } from "react";
import { Link, useNavigate, useParams } from "react-router-dom";
import {
  Accordion,
  AccordionDetails,
  AccordionSummary,
  Alert,
  Box,
  Breadcrumbs,
  Button,
  Container,
  Divider,
  IconButton,
  Paper,
  TextField,
  Typography,
} from "@mui/material";
import ExpandMoreIcon from "@mui/icons-material/ExpandMore";
import Grid from "@mui/joy/Grid";
import DeleteIcon from "@mui/icons-material/Delete";
import VisibilityIcon from "@mui/icons-material/Visibility";
import VisibilityOffIcon from "@mui/icons-material/VisibilityOff";
import ReactQuill from "react-quill";
import "react-quill/dist/quill.snow.css";

import { getVersion, postVersion } from "../api/VersionApi.js";
import { getEntry, postEntry } from "../api/EntryApi.js";
import { getWiki } from "../api/WikiApi.js";
import { getMedia, postMedia } from "../api/MediaApi.js";
import ConfirmationModal from "../components/ConfirmationModal.jsx";
import theme from "../styles/theme.js";

function FormEntradaPage() {
  // Destructure the URL parameters
  const { entryId, versionId, wikiId } = useParams();
  const navigate = useNavigate();

  // Define state variables
  const [entry, setEntry] = useState({ title: "", author: "" });
  const [wiki, setWiki] = useState({});
  const [version, setVersion] = useState({ content: "", address: "" });
  const [formErrors, setFormErrors] = useState({});
  const [uploads, setUploads] = useState([]);
  const [existingImages, setExistingImages] = useState([]);
  const [error, setError] = useState(null);
  const [openConfirmDialog, setOpenConfirmDialog] = useState(false);
  const [versionError, setVersionError] = useState(null);

  const formRef = useRef(null);
  const isLoggedIn = !!sessionStorage.getItem("appUser");
  const userId = sessionStorage.getItem("id");

  // Determine if we're creating a new entry
  const isNewEntry = !entryId;

  // File upload settings
  const MAX_FILE_SIZE = 5 * 1024 * 1024; // 5MB
  const ALLOWED_TYPES = ["image/jpeg", "image/png", "image/gif"];

  // Redirect if not logged in
  useEffect(() => {
    if (!isLoggedIn) {
      navigate("/");
    }
  }, [isLoggedIn, navigate]);

  // Fetch the wiki details
  useEffect(() => {
    const id = wikiId || entry.wiki_id;
    if (id) {
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
  }, [entry, wikiId]);

  // Fetch the entry details if editing
  useEffect(() => {
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
    }
  }, [entryId]);

  // Fetch the version details if editing
  useEffect(() => {
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
  }, [versionId]);

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

  // Handler for ReactQuill editor change
  const handleEditorChange = (content) => {
    setVersion((prevVersion) => ({
      ...prevVersion,
      content: content,
    }));
    setFormErrors((prevErrors) => ({
      ...prevErrors,
      content: "",
    }));
  };

  // Handler for TextField change in version data
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

  // Handler for TextField change in entry data
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

  // Toggle visibility of existing images
  const handleToggleVisibility = (index) => {
    setExistingImages((prevImages) =>
      prevImages.map((image, i) =>
        i === index ? { ...image, isVisible: !image.isVisible } : image
      )
    );
  };

  // Validation function
  const validate = () => {
    let isValid = true;
    const errors = {};

    if (
      !version.content ||
      version.content.replace(/<[^>]+>/g, "").trim().length === 0
    ) {
      errors.content = "Introduzca contenido";
      isValid = false;
    }
    if (!entry.title || entry.title.trim() === "") {
      errors.title = "Introduzca un título";
      isValid = false;
    }

    setFormErrors(errors);
    return isValid;
  };

  // Handle form submission
  const handleSubmit = (event) => {
    event.preventDefault();
    if (validate()) {
      setOpenConfirmDialog(true);
    }
  };

  const handleConfirmSubmit = async () => {
    setOpenConfirmDialog(false);

    if (!validate()) {
      return;
    }

    let newMediaIds = [];

    const entryData = {
      title: entry.title,
      wiki_id: wiki.id,
      author: userId,
    };

    // Upload new images
    if (uploads.length > 0) {
      for (const file of uploads) {
        const formData = new FormData();
        formData.append("image", file);

        try {
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
      editor: userId,
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
      } else {
        versionData.entry_id = entryId;
        // Update the entry if necessary
        await postEntry({ ...entryData, id: entryId });
      }

      await postVersion(versionData);
      navigate(`/entrada/${versionData.entry_id}`);
    } catch (error) {
      console.error("Error posting version:", error);
    }
  };

  const handleCancelSubmit = () => {
    setOpenConfirmDialog(false);
  };

  if (!isLoggedIn) {
    return null;
  }

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
          <Grid container spacing={2}>
            <Grid xs={12}>
              <Typography variant="h4" gutterBottom>
                {isNewEntry ? "Crear Entrada" : "Editar Entrada"}
              </Typography>
            </Grid>
            <Grid xs={12} sm={6}>
              <TextField
                id="title"
                name="title"
                label="Título"
                value={entry.title || ""}
                onChange={handleEntryChange}
                variant="outlined"
                fullWidth
                error={!!formErrors.title}
                helperText={formErrors.title}
                slotProps={{
                  inputLabel: {
                    shrink: true,
                  },
                }}
              />
            </Grid>
            <Grid xs={12}>
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
          {versionError && (
            <Typography variant="body1" color="error" gutterBottom>
              {versionError}
            </Typography>
          )}
          <Typography variant="subtitle1" sx={{ mt: 3, mb: 1 }}>
            Contenido:
          </Typography>
          <Box sx={{ height: "400px" }}>
            <ReactQuill
              theme="snow"
              value={version.content || ""}
              onChange={handleEditorChange}
              style={{
                height: "100%",
                display: "flex",
                flexDirection: "column",
                border: formErrors.content
                  ? `1px solid ${theme.palette.error.main}`
                  : "1px solid #ccc",
              }}
            />
            {formErrors.content && (
              <Typography
                variant="body2"
                sx={{
                  color: theme.palette.error.main,
                  fontSize: "12px",
                  mt: 1,
                }}
              >
                {formErrors.content}
              </Typography>
            )}
          </Box>
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
                      {image.isVisible
                        ? <VisibilityIcon />
                        : <VisibilityOffIcon />}
                    </IconButton>
                    <IconButton
                      onClick={() => handleRemoveImage(index, true)}
                      sx={{ ml: 1 }}
                    >
                      <DeleteIcon />
                    </IconButton>
                  </Box>
                ))}
              </AccordionDetails>
            </Accordion>
          )}
          {error && (
            <Alert severity="error" sx={{ mt: 2 }}>
              {error}
            </Alert>
          )}
          <Box sx={{ mt: 5, pt: 1 }} display="flex" justifyContent="flex-end">
            <Button
              type="submit"
              variant="contained"
              color="primary"
              sx={{ mr: 2 }}
            >
              Enviar
            </Button>
            <Button
              variant="outlined"
              color="secondary"
              onClick={() => navigate(-1)}
            >
              Cancelar
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

export default FormEntradaPage;
