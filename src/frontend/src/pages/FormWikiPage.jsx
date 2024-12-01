import { useEffect, useState } from "react";
import {
  Container,
  Typography,
  TextField,
  Button,
  Box,
  IconButton,
  Alert,
  Breadcrumbs
} from "@mui/material";
import { useNavigate, useParams, Link } from "react-router-dom";
import { getWiki, postWiki, putWiki } from "../api/WikiApi";
import { postMedia } from "../api/MediaApi";
import DeleteIcon from "@mui/icons-material/Delete";
import ConfirmationModal from "../components/ConfirmationModal.jsx";
import { useToast } from "../context/ToastContext.jsx";

function FormWikiPage() {
  const { wikiId } = useParams();
  const MAX_FILE_SIZE = 5 * 1024 * 1024; // 5MB
  const ALLOWED_TYPES = ["image/jpeg", "image/png", "image/gif"];
  const [wiki, setWiki] = useState({
    title: "",
    description: "",
    category: "",
    media_id: ""
  });
  const [error, setError] = useState(null);
  const navigate = useNavigate();

  const { showToast } = useToast();
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [formErrors, setFormErrors] = useState({
    title: "",
    description: "",
    category: "",
    media_id: ""
  });

  const [upload, setUpload] = useState(null);
  const [uploading, setUploading] = useState(false);

  useEffect(() => {
    if (wikiId) {
      getWiki(wikiId)
        .then((data) => {
          if (data && Object.keys(data).length > 0) {
            setWiki(data);
            if (data.media_id && data.media_id.length > 0) {
              // Fetch and set the uploaded image if needed
            }
          } else {
            setError("Wiki no encontrada.");
          }
        })
        .catch(() => setError("Error al obtener los detalles de la wiki."));
    }
  }, [wikiId]);

  const handleImageChange = async (event) => {
    const file = event.target.files[0];
    setError(null);

    if (!file) return;

    if (!ALLOWED_TYPES.includes(file.type)) {
      setError(`Tipo de archivo no permitido: ${file.name}`);
      return;
    }

    if (file.size > MAX_FILE_SIZE) {
      setError(`Archivo demasiado grande: ${file.name}`);
      return;
    }

    setUploading(true);

    const formData = new FormData();
    formData.append("image", file);

    try {
      const response = await postMedia(formData);
      setUpload({ file, id: response.id });
    } catch (error) {
      console.error(`Error uploading file ${file.name}:`, error);
      setError(`Error uploading file ${file.name}`);
    } finally {
      setUploading(false);
    }
  };

  const handleRemoveImage = () => {
    setUpload(null);
  };

  const handleChange = (event) => {
    const { name, value } = event.target;
    setWiki((prevWiki) => ({
      ...prevWiki,
      [name]: value,
    }));
    // Clear error for the field
    setFormErrors((prevErrors) => ({
      ...prevErrors,
      [name]: "",
    }));
  };

  const validate = () => {
    let isValid = true;
    let errors = { title: "", description: "", category: "" };

    if (!wiki.title.trim()) {
      errors.title = "Introduzca un título";
      isValid = false;
    }
    if (!wiki.description.trim()) {
      errors.description = "Introduzca una descripción";
      isValid = false;
    }
    if (!wiki.category.trim()) {
      errors.category = "Introduzca una categoría";
      isValid = false;
    }

    setFormErrors(errors);
    return isValid;
  };

  const handleSubmit = async () => {
    try {
      const mediaID = upload.id ? upload.id : null;
      const wikiData = {
        ...wiki,
        media_id: mediaID ? mediaID : null,
      };

      if (wikiId) {
        await putWiki(wikiId, wikiData);
        showToast("Wiki actualizada correctamente", "success");
      } else {
        await postWiki(wikiData);
        showToast("Wiki creada correctamente", "success");
      }
      navigate("/");
    } catch (error) {
      console.error("Error al guardar la wiki:", error);
      setError("Error al guardar la wiki.");
    }
  };

  const onSubmit = (event) => {
    event.preventDefault();
    if (validate()) {
      setIsModalOpen(true);
    }
  };

  return (
    <Container maxWidth="md" sx={{ mt: 4, mb: 4 }}>
      <Breadcrumbs aria-label="breadcrumb" sx={{ mb: 2 }}>
        <Typography className="breadcrumb-link" component={Link} to="/">
          Inicio
        </Typography>
        <Typography className="breadcrumb-link" component={Link} to={`/wiki/${wikiId}`}>
          {wiki.title}
        </Typography>
      </Breadcrumbs>

      <Typography variant="h4" gutterBottom>
        {wikiId ? "Editar Wiki" : "Crear Nueva Wiki"}
      </Typography>
      {(error || Object.values(formErrors).some((e) => e)) && (
        <Alert severity="error" sx={{ mb: 2 }}>
          {error ||
            Object.values(formErrors)
              .filter((e) => e)
              .join(", ")}
        </Alert>
      )}
      <Box component="form" onSubmit={onSubmit} noValidate sx={{ mt: 1 }}>
        <TextField
          margin="normal"
          required
          fullWidth
          id="title"
          name="title"
          label="Título"
          value={wiki.title}
          onChange={handleChange}
          variant="outlined"
          error={!!formErrors.title}
          helperText={formErrors.title}
        />
        <TextField
          margin="normal"
          required
          fullWidth
          id="description"
          name="description"
          label="Descripción"
          value={wiki.description}
          onChange={handleChange}
          variant="outlined"
          multiline
          rows={4}
          error={!!formErrors.description}
          helperText={formErrors.description}
        />
        <TextField
          margin="normal"
          required
          fullWidth
          id="category"
          name="category"
          label="Categoría"
          value={wiki.category}
          onChange={handleChange}
          variant="outlined"
          error={!!formErrors.category}
          helperText={formErrors.category}
        />
        <Button variant="contained" component="label" sx={{ mt: 2 }}>
          Añadir Imagen
          <input
            id="image-input"
            type="file"
            hidden
            accept="image/*"
            onChange={handleImageChange}
          />
        </Button>
        {upload && (
          <Box sx={{ mt: 2, display: "flex", alignItems: "center" }}>
            <Typography variant="body2">{upload.file.name}</Typography>
            <IconButton onClick={handleRemoveImage} sx={{ ml: 1 }}>
              <DeleteIcon />
            </IconButton>
          </Box>
        )}
        <Button
          type="submit"
          fullWidth
          variant="contained"
          sx={{ mt: 3, mb: 2 }}
          disabled={uploading}
        >
          {uploading
            ? "Subiendo..."
            : wikiId
            ? "Guardar Cambios"
            : "Crear Wiki"}
        </Button>
      </Box>

      <ConfirmationModal
        show={isModalOpen}
        handleClose={() => setIsModalOpen(false)}
        handleConfirm={handleSubmit}
        message={`¿Estás seguro de que quieres ${
          wikiId ? "guardar los cambios" : "crear esta wiki"
        }?`}
      />
    </Container>
  );
}

export default FormWikiPage;
