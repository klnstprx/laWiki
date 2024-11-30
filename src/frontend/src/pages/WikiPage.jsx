import { useEffect, useState } from "react";
import { useParams, Link } from "react-router-dom";
import {
  Container,
  Paper,
  Typography,
  Button,
  Alert,
  List,
  ListItem,
} from "@mui/material";
import { searchEntries } from "../api/EntryApi.js";
import { getWiki } from "../api/WikiApi.js";
import EntradaCard from "../components/EntradaCard.jsx";

function WikiPage() {
  const [wiki, setWiki] = useState({});
  const [entradas, setEntradas] = useState([]);
  const [error, setError] = useState(null);

  const { id } = useParams();

  useEffect(() => {
    getWiki(id)
      .then((data) => {
        if (data && Object.keys(data).length > 0) {
          setWiki(data);
        } else {
          setError("Wiki no encontrada.");
        }
      })
      .catch((err) => setError(err.message));
  }, [id]);

  useEffect(() => {
    searchEntries({ wikiID: id })
      .then((data) => {
        if (data && Array.isArray(data)) {
          setEntradas(data);
        } else {
          setEntradas([]); // Ensure entradas is always an array
        }
      })
      .catch((err) => setError(err.message));
  }, [id]);

  return (
    <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
      {error && (
        <Alert severity="error" sx={{ mb: 2 }}>
          {error}
        </Alert>
      )}

      {!error && wiki && Object.keys(wiki).length > 0 && (
        <>
          {/* Page Header */}
          <Paper
            elevation={3}
            sx={{ p: 2, textAlign: "center", borderRadius: 1 }}
          >
            <Typography variant="h3" component="h1" sx={{ m: 0 }}>
              {wiki.title}
            </Typography>
          </Paper>

          {/* Wiki Information */}
          <Paper elevation={3} sx={{ p: 3, my: 4, borderRadius: 1 }}>
            <Typography variant="h6" gutterBottom>
              <strong>Título:</strong> {wiki.title}
            </Typography>
            <Typography variant="h6" gutterBottom>
              <strong>Descripción:</strong> {wiki.description}
            </Typography>
            <Typography variant="h6" gutterBottom>
              <strong>Categoría:</strong> {wiki.category}
            </Typography>
          </Paper>

          {/* Entradas */}
          <Paper elevation={3} sx={{ p: 3, mb: 4, borderRadius: 1 }}>
            <Typography
              variant="h4"
              component="h2"
              sx={{ borderBottom: "1px solid", pb: 1, mb: 2 }}
            >
              Entradas
            </Typography>
            {entradas && entradas.length > 0 ? ( // Added null check for entradas
              <List>
                {entradas.map((entrada) => (
                  <ListItem key={entrada.id} divider>
                    <EntradaCard
                      id={entrada.id}
                      title={entrada.title}
                      author={entrada.author}
                      createdAt={entrada.created_at}
                    />
                  </ListItem>
                ))}
              </List>
            ) : (
              <Alert severity="info">No entries found.</Alert>
            )}
          </Paper>

          {/* Button to create new entry */}
          <Button
            component={Link}
            to={`/crear-entrada/${id}`}
            variant="contained"
            color="primary"
            sx={{ mt: 2 }}
          >
            Crear Nueva Entrada
          </Button>
        </>
      )}
    </Container>
  );
}

export default WikiPage;
