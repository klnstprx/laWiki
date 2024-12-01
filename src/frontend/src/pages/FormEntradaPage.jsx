import { useState } from "react";
import {
  Container,
  Typography,
  TextField,
  Button,
  Box,
  IconButton,
  Alert,
} from "@mui/material";
import { useNavigate, useParams, Link } from "react-router-dom";
import { postEntry } from "../api/EntryApi";
import { postMedia } from "../api/MediaApi";
import DeleteIcon from "@mui/icons-material/Delete";
import { useToast } from "../context/ToastContext.jsx";
import { Breadcrumbs } from "@mui/material";
import { useEffect } from "react";
import { getWiki } from "../api/WikiApi";

function FormEntradaPage() {
  const { id: wikiId } = useParams();
  const [wiki, setWiki] = useState({});
  const [title, setTitle] = useState("");
  const [uploads, setUploads] = useState([]);
  const navigate = useNavigate();
  const { showToast } = useToast();
  const [error, setError] = useState(null);
  const [titleError, setTitleError] = useState("");
  const [uploading, setUploading] = useState(false);
  const [loading, setLoading] = useState(false);

  const MAX_FILE_SIZE = 5 * 1024 * 1024; // 5MB
  const ALLOWED_TYPES = ["image/jpeg", "image/png", "image/gif"];

  const handleImageChange = async (event) => {
    const files = Array.from(event.target.files);
    setError(null);
    setUploading(true);

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
      setUploading(false);
      return;
    }

    const newUploads = [];

    for (const file of validFiles) {
      const formData = new FormData();
      formData.append("image", file);

      try {
        const response = await postMedia(formData);
        newUploads.push({ file, id: response.id });
      } catch (error) {
        console.error(`Error uploading file ${file.name}:`, error);
        setError(`Error uploading file ${file.name}`);
      }
    }

    setUploads((prevUploads) => [...prevUploads, ...newUploads]);
    setUploading(false);
  };

  const handleRemoveImage = (index) => {
    setUploads((prevUploads) => prevUploads.filter((_, i) => i !== index));
  };

  const handleSubmit = async (event) => {
    event.preventDefault();
    setError(null);

    let isValid = true;

    if (!title.trim()) {
      setTitleError("Introduzca un título");
      isValid = false;
    } else {
      setTitleError("");
    }

    if (!isValid) return;

    const mediaIDs = uploads.map((upload) => upload.id);

    const entryData = {
      title,
      wiki_id: wikiId,
      media_ids: mediaIDs,
    };

    setLoading(true);

    try {
      await postEntry(entryData);
      navigate(`/wiki/${wikiId}`);
      showToast("Entrada creada correctamente", "success");
    } catch (error) {
      const errorMessage =
        error.response?.data?.message || "Error al crear la entrada";
      setError(errorMessage);
      console.error("Error al crear la entrada:", error);
    } finally {
      setLoading(false);
    }
  };

  // Fetch the wiki details
  useEffect(() => {
    if (wikiId) {
      getWiki(wikiId)
        .then((data) => {
          if (data && Object.keys(data).length > 0) {
            setWiki(data);
          } else {
            setError("No se encontró la wiki asociada a esta entrada.");
          }
        })
        .catch(() =>
          setError("Se produjo un error al obtener la wiki asociada."),
        );
    }
  }, [wikiId]);

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
        Crear Nueva Entrada
      </Typography>
      {(error) && (
        <Alert severity="error" sx={{ mb: 2 }}>
          {error}
        </Alert>
      )}
      <Box component="form" onSubmit={handleSubmit} sx={{ mt: 1 }}>
        <TextField
          margin="normal"
          fullWidth
          id="title"
          label="Título"
          name="title"
          autoComplete="off"
          autoFocus
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          error={!!titleError}
          helperText={titleError}
        />
        <Button variant="contained" component="label" sx={{ mt: 2 }}>
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
          <Box sx={{ mt: 2 }}>
            {uploads.map((upload, index) => (
              <Box
                key={upload.id || index}
                sx={{ display: "flex", alignItems: "center", mt: 1 }}
              >
                <Typography variant="body2">{upload.file.name}</Typography>
                <IconButton
                  onClick={() => handleRemoveImage(index)}
                  sx={{ ml: 1 }}
                >
                  <DeleteIcon />
                </IconButton>
              </Box>
            ))}
          </Box>
        )}
        <Box sx={{ display: "flex", gap: 2, mt: 2 }}>
          <Button
            type="submit"
            variant="contained"
            color="primary"
            disabled={uploading || loading}
          >
            {uploading || loading ? "Creando..." : "Crear Entrada"}
          </Button>
          <Button
            component={Link}
            to={`/wiki/${wikiId}`}
            variant="outlined"
            color="primary"
          >
            Cancelar
          </Button>
        </Box>
      </Box>
    </Container>
  );
}

export default FormEntradaPage;
