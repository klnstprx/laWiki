import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import {
  Container,
  Paper,
  Typography,
  Alert,
  List,
  ListItem,
} from "@mui/material";
import VersionCard from "../components/VersionCard.jsx";
import { searchVersions } from "../api/VersionApi.js";

function VersionPage() {
  const [versiones, setVersiones] = useState([]);
  const [error, setError] = useState(null);

  const { entryId } = useParams();

  useEffect(() => {
    if (entryId) {
      searchVersions({ entryID: entryId })
        .then((data) => {
          if (data && data.length > 0) {
            setVersiones(data);
          } else {
            setError("No se encontraron versiones para esta entrada.");
          }
        })
        .catch(() =>
          setError("Se ha producido un error al obtener las versiones."),
        );
    } else {
      setError("No se proporcionó un ID de entrada válido.");
    }
  }, [entryId]);

  return (
    <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
      <Paper elevation={3} sx={{ p: 3, textAlign: "center", mb: 4 }}>
        <Typography variant="h3" component="h1">
          Versiones de la Entrada
        </Typography>
        {error && (
          <Alert severity="error" sx={{ mt: 2 }}>
            {error}
          </Alert>
        )}
      </Paper>

      {!error && versiones.length > 0 ? (
        <Paper elevation={3} sx={{ p: 3 }}>
          <List>
            {versiones.map((version) => (
              <ListItem key={version.id} divider>
                <VersionCard
                  entradaId={entryId}
                  versionId={version.id}
                  editor={version.editor}
                  created_at={version.created_at}
                />
              </ListItem>
            ))}
          </List>
        </Paper>
      ) : (
        !error && (
          <Alert severity="info" sx={{ mt: 2 }}>
            No se encontraron versiones para esta entrada.
          </Alert>
        )
      )}
    </Container>
  );
}

export default VersionPage;
