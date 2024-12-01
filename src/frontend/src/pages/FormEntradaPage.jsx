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
import DeleteIcon from "@mui/icons-material/Delete";
import { useToast } from "../context/ToastContext.jsx";

function FormEntradaPage() {
  const { id: wikiId } = useParams();
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [image, setImage] = useState(null);
  const navigate = useNavigate();
  const { showToast } = useToast();
  const [error, setError] = useState(null);

  const handleImageChange = (event) => {
    setImage(event.target.files[0]);
  };

  const handleSubmit = async (event) => {
    event.preventDefault();
    const form = event.target;
    const formData = new FormData(form);
    const jsonData = {};
    formData.forEach((value, key) => {
      jsonData[key] = value;
    });

    jsonData["wiki_id"] = wikiId;

    try {
      await postEntry(jsonData);
      navigate(`/wiki/${wikiId}`);
      showToast("Entrada creada correctamente", "success");
    } catch (error) {
      setError("Error al crear la entrada");
      console.error("Error al crear la entrada:", error);
    }
  };

  return (
    <Container maxWidth="md">
      <Typography variant="h4" gutterBottom>
        Crear Nueva Entrada
      </Typography>
      {error && (
        <Alert severity="error" sx={{ mb: 2 }}>
          {error}
        </Alert>
      )}
      <Box component="form" onSubmit={handleSubmit} noValidate sx={{ mt: 1 }}>
        <TextField
          margin="normal"
          required
          fullWidth
          id="title"
          label="Título"
          name="title"
          autoComplete="title"
          autoFocus
          value={title}
          onChange={(e) => setTitle(e.target.value)}
        />
        <TextField
          margin="normal"
          required
          fullWidth
          id="description"
          label="Descripción"
          name="description"
          autoComplete="description"
          value={description}
          onChange={(e) => setDescription(e.target.value)}
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
        {image && (
          <Box sx={{ display: "flex", alignItems: "center", mt: 2 }}>
            <Typography variant="body2">{image.name}</Typography>
            <IconButton onClick={() => setImage(null)} sx={{ ml: 1 }}>
              <DeleteIcon />
            </IconButton>
          </Box>
        )}
        <Button type="submit" variant="contained" color="primary" fullWidth>
          Crear Entrada
        </Button>
        <Button
          component={Link}
          to={`/wiki/${wikiId}`}
          variant="outlined"
          color="primary"
          fullWidth
        >
          Cancelar
        </Button>
      </Box>
    </Container>
  );
}

export default FormEntradaPage;
