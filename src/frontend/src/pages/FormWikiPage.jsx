import { useEffect, useState, useRef } from "react";
import { getWiki, postWiki, putWiki } from "../api/WikiApi.js";
import { useNavigate, useParams } from "react-router-dom";
import {
  Container,
  Paper,
  Typography,
  Button,
  TextField,
  Grid2,
  Box,
} from "@mui/material";

function FormWikiPage() {
  const { wikiId } = useParams();
  const [wiki, setWiki] = useState({
    title: "",
    description: "",
    category: "",
  });
  const [error, setError] = useState(null);
  const formRef = useRef(null);
  const navigate = useNavigate();

  // Obtener detalles de la wiki si se está editando
  useEffect(() => {
    if (wikiId) {
      getWiki(wikiId)
        .then((data) => {
          if (data && Object.keys(data).length > 0) {
            setWiki(data);
          } else {
            setError("No se encontró la wiki solicitada.");
          }
        })
        .catch(() =>
          setError("Se produjo un error al obtener los detalles de la wiki.")
        );
    }
  }, [wikiId]);

  // Manejador de cambios en los campos del formulario
  const handleChange = (event) => {
    const { name, value } = event.target;
    setWiki((prevWiki) => ({
      ...prevWiki,
      [name]: value,
    }));
  };

  // Manejador para enviar el formulario
  async function handleSubmit(event) {
    event.preventDefault();

    const wikiData = {
      title: wiki.title,
      description: wiki.description,
      category: wiki.category,
    };

    try {
      if (wikiId) {
        // Actualizar wiki existente
        await putWiki(wikiId, wikiData);
        navigate(`/wiki/${wikiId}`);
      } else {
        // Crear nueva wiki
        const newWiki = await postWiki(wikiData);
        navigate(`/wiki/${newWiki.id}`);
      }
    } catch (error) {
      console.error("Error al guardar la wiki:", error);
      setError("Error al guardar la wiki.");
    }
  }

  return (
    <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
      <Paper elevation={3} sx={{ p: { xs: 2, md: 4 }, mb: 4 }}>
        <form id="formWiki" ref={formRef} onSubmit={handleSubmit}>
          <Typography variant="h4" gutterBottom>
            {wikiId ? "Editar Wiki" : "Crear Nueva Wiki"}
          </Typography>
          {error && (
            <Typography variant="body1" color="error" gutterBottom>
              {error}
            </Typography>
          )}
          <Grid2 container spacing={2}>
            <Grid2 xs={12}>
              <TextField
                id="title"
                name="title"
                label="Título"
                value={wiki.title}
                onChange={handleChange}
                variant="outlined"
                fullWidth
                required
              />
            </Grid2>
            <Grid2 xs={12}>
              <TextField
                id="description"
                name="description"
                label="Descripción"
                value={wiki.description}
                onChange={handleChange}
                variant="outlined"
                fullWidth
                multiline
                rows={4}
                required
              />
            </Grid2>
            <Grid2 xs={12}>
              <TextField
                id="category"
                name="category"
                label="Categoría"
                value={wiki.category}
                onChange={handleChange}
                variant="outlined"
                fullWidth
                required
              />
            </Grid2>
          </Grid2>
          <Box sx={{ mt: 5 }} display="flex" justifyContent="flex-end">
            <Button type="submit" variant="contained" color="primary">
              {wikiId ? "Guardar Cambios" : "Crear Wiki"}
            </Button>
          </Box>
        </form>
      </Paper>
    </Container>
  );
}

export default FormWikiPage;