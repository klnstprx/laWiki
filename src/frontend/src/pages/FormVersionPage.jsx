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
import { getEntry } from "../api/EntryApi.js";
import { getWiki } from "../api/WikiApi.js";
import { postMedia, getMedia} from "../api/MediaApi";
import DeleteIcon from "@mui/icons-material/Delete";
import VisibilityIcon from '@mui/icons-material/Visibility';
import VisibilityOffIcon from '@mui/icons-material/VisibilityOff';
import ExpandMoreIcon from "@mui/icons-material/ExpandMore";

function FormVersionPage() {
  const { entryId, versionId } = useParams();
  const [entry, setEntry] = useState({});
  const [wiki, setWiki] = useState({});
  const [version, setVersion] = useState({});
  const [versionError, setVersionError] = useState(null);
  const [formErrors, setFormErrors] = useState({
    editor: "",
    content: "",
  });
  const [uploads, setUploads] = useState([]); // Store selected files
  const [existingImages, setExistingImages] = useState([]);
  const [error, setError] = useState(null);
  const formRef = useRef(null);
  const navigate = useNavigate();

  const MAX_FILE_SIZE = 5 * 1024 * 1024; // 5MB
  const ALLOWED_TYPES = ["image/jpeg", "image/png", "image/gif"];
  const isLoggedIn = !!sessionStorage.getItem('user'); // Suponiendo que guardas el estado de inicio de sesión en sessionStorage
  const id = sessionStorage.getItem('id'); // Suponiendo que guardas el id del usuario en sessionStorage

  useEffect(() => {
    if (!isLoggedIn) {
      navigate('/login');
    }
  }, [isLoggedIn, navigate]);

  // Fetch the version details
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

  // Fetch the entry details
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
    } else {
      setVersionError("No se proporcionó un ID de entrada válido.");
    }
  }, [entryId]);

  // Fetch the wiki details
  useEffect(() => {
    if (entry && entry.wiki_id) {
      getWiki(entry.wiki_id)
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
  }, [entry, entryId]);

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
  const handleChange = (event) => {
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
    if (
      !version.content ||
      version.content.replace(/<[^>]+>/g, "").trim().length === 0
    ) {
      errors.content = "Introduzca contenido";
      isValid = false;
    }
    setFormErrors(errors);
    return isValid;
  };

  // Handler to submit the version
  const handleSubmit = async (event) => {
    event.preventDefault();

    if (!validate()) {
      return;
    }

    let newMediaIds = [];

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

    const jsonData = {
      content: version.content,
      editor: id,
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
      await postVersion(jsonData);
      navigate(`/entrada/${entryId}`);
    } catch (error) {
      console.error("Error posting version:", error);
    }
  };

  if (!isLoggedIn) {
    return null; // O puedes mostrar un mensaje de carga o redirección
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
        <Typography
          className="breadcrumb-link"
          component={Link}
          to={`/entrada/${entry.id}`}
        >
          {entry.title}
        </Typography>
      </Breadcrumbs>

      <Paper elevation={3} sx={{ p: { xs: 2, md: 4 }, mb: 4 }}>
        <form id="miFormulario" ref={formRef} onSubmit={handleSubmit}>
          <Grid container spacing={2} alignItems="center">
            <Grid xs={12} sm={8}>
              <Typography variant="h4">Editar entrada</Typography>
            </Grid>
            <Grid xs={12} sm={4}>
              <TextField
                id="address"
                name="address"
                label="Ubicación (opcional)"
                value={version.address || ""}
                onChange={handleChange}
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
                    <a href={image.uploadUrl} target="_blank" rel="noopener noreferrer">
                      <Typography variant="body2" color="primary" sx={{ textDecoration: 'underline', cursor: 'pointer' }}>
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
    </Container>
  );
}

export default FormVersionPage;