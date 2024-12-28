import { useState } from "react";
import {
  Alert,
  Box,
  Button,
  Container,
  TextField,
  Typography,
} from "@mui/material";
import { Link, useNavigate, useParams } from "react-router-dom";
import { postEntry } from "../api/EntryApi";
import { useToast } from "../context/ToastContext.jsx";
import { Breadcrumbs } from "@mui/material";
import { useEffect } from "react";
import { getWiki } from "../api/WikiApi";
import { useAuth } from "../context/AuthContext";

function FormEntradaPage() {
  const { id: wikiId } = useParams();
  const [wiki, setWiki] = useState({});
  const [title, setTitle] = useState("");
  const navigate = useNavigate();
  const { showToast } = useToast();
  const [error, setError] = useState(null);
  const [titleError, setTitleError] = useState("");
  const { user } = useAuth();
  const isLoggedIn = !!user;

  const handleSubmit = async (event) => {
    event.preventDefault();
    const idAutor = user.id;
    setError(null);

    let isValid = true;

    if (!title.trim()) {
      setTitleError("Introduzca un título");
      isValid = false;
    } else {
      setTitleError("");
    }

    if (!isValid) return;

    const entryData = {
      title,
      wiki_id: wikiId,
      author: idAutor,
    };

    try {
      await postEntry(entryData);
      navigate(`/wiki/${wikiId}`);
      showToast("Entrada creada correctamente", "success");
    } catch (error) {
      const errorMessage = error.response?.data?.message ||
        "Error al crear la entrada";
      setError(errorMessage);
      console.error("Error al crear la entrada:", error);
    }
  };

  useEffect(() => {
    if (!isLoggedIn) {
      navigate("/");
    }
  }, [isLoggedIn, navigate]);

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
          setError("Se produjo un error al obtener la wiki asociada.")
        );
    }
  }, [wikiId]);

  return (
    <Container maxWidth="md" sx={{ mt: 4, mb: 4 }}>
      <Breadcrumbs aria-label="breadcrumb" sx={{ mb: 2 }}>
        <Typography className="breadcrumb-link" component={Link} to="/">
          Inicio
        </Typography>
        <Typography
          className="breadcrumb-link"
          component={Link}
          to={`/wiki/${wikiId}`}
        >
          {wiki.title}
        </Typography>
      </Breadcrumbs>

      <Typography variant="h4" gutterBottom>
        Crear Nueva Entrada
      </Typography>
      {error && (
        <Alert severity="error" sx={{ mb: 2 }}>
          {error}
        </Alert>
      )}
      {isLoggedIn && (
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
            onChange={(e) =>
              setTitle(e.target.value)}
            error={!!titleError}
            helperText={titleError}
            slotProps={{
              inputLabel: {
                shrink: true,
              },
            }}
          />
          <Box sx={{ display: "flex", gap: 2, mt: 2 }}>
            <Button
              type="submit"
              variant="contained"
              color="primary"
            >
              {"Crear Entrada"}
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
      )}
      {!isLoggedIn && <Alert sx={{ mb: 2 }}>No estas logeado.</Alert>}
    </Container>
  );
}

export default FormEntradaPage;
