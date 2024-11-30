import { useState } from "react";
import { postEntry } from "../api/EntryApi.js";
import {
  Container,
  Paper,
  Typography,
  Button,
  TextField,
  Alert,
  Grid2,
} from "@mui/material";
import { useNavigate, useParams, Link } from "react-router-dom";
import { useToast } from "../context/ToastContext.jsx";

function PostEntradaPage() {
  const [error, setError] = useState(null);
  const navigate = useNavigate();
  const { showToast } = useToast();
  const { id } = useParams();

  async function enviarJSON(event) {
    event.preventDefault();

    const form = event.target;
    const formData = new FormData(form);

    const jsonData = {};
    formData.forEach((value, key) => {
      jsonData[key] = value;
    });

    // Añadir el wiki_id automáticamente
    jsonData["wiki_id"] = id;

    try {
      const result = await postEntry(jsonData);
      console.log("Respuesta del servidor:", result);

      showToast("Entrada creada correctamente", "success");
      navigate(`/wiki/${id}`);
    } catch (error) {
      setError("Error al crear la entrada");
      console.error("Error al enviar:", error);
    }
  }

  return (
    <Container maxWidth="sm" sx={{ mt: 4 }}>
      <Paper elevation={3} sx={{ p: 4 }}>
        <Typography variant="h4" component="h2" align="center" gutterBottom>
          Crear Nueva Entrada
        </Typography>

        {error && (
          <Alert severity="error" sx={{ mb: 2 }}>
            {error}
          </Alert>
        )}

        <form onSubmit={enviarJSON}>
          <Grid2 container spacing={2}>
            <Grid2 xs={12}>
              <TextField label="Título" name="title" required fullWidth />
            </Grid2>

            <Grid2 xs={12}>
              <TextField label="Autor" name="author" required fullWidth />
            </Grid2>

            <Grid2 xs={12}>
              <Button
                type="submit"
                variant="contained"
                color="primary"
                fullWidth
              >
                Crear Entrada
              </Button>
            </Grid2>
            <Grid2 xs={12}>
              <Button
                component={Link}
                to={`/wiki/${id}`}
                variant="outlined"
                color="primary"
                fullWidth
              >
                Cancelar
              </Button>
            </Grid2>
          </Grid2>
        </form>
      </Paper>
    </Container>
  );
}

export default PostEntradaPage;
