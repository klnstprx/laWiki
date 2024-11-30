import { useEffect, useState, useRef } from "react";
import {
  postVersion,
  getVersion
} from "../api/VersionApi.js";
import { getEntry } from "../api/EntryApi.js";
import { useSearchParams, useParams } from "react-router-dom";
import Comentario from "../components/Comentario.jsx";
import Version from "../components/Version.jsx";
import MainLayout from "../layout/MainLayout.jsx";
import ConfirmationModal from "../components/ConfirmationModal.jsx";
import { useToast } from "../context/ToastContext.1.jsx";
import {
  Container,
  Paper,
  Typography,
  Button,
  Alert,
  List,
  ListItem,
  Input,
  Grid2,
} from "@mui/material";

function EditarEntradaPage() {
  const { entryId, versionId } = useParams();
  const [entrada, setEntrada] = useState({});
  const [version, setVersion] = useState({});
  const [versionError, setVersionError] = useState(null);

  const formRef = useRef(null);

  // Obtener la versión
  useEffect(() => {
    if (versionId) {
      getVersion(versionId)
        .then((data) => {
          if (data && Object.keys(data).length > 0) {
            setVersion(data);
          } else {
            setVersionError("No se encontró la versión solicitada.");
          }
        })
        .catch(() =>
          setVersionError("Se produjo un error al obtener la versión."),
        );
    } else {
      setVersionError("No se proporcionó un ID de versión válido.");
    }
  }, [versionId]);

  const handleChange = (event) => {
    const { name, value } = event.target;
    setVersion((prevVersion) => ({
      ...prevVersion,
      [name]: value,
    }));
  };

  // Handler para enviar el comentario
  async function subirVersion(event) {
    event.preventDefault();

    const formData = new FormData(event.target);
    const jsonData = Object.fromEntries(formData.entries());

    jsonData["entry_id"] = entryId;

    postVersion(jsonData);

    window.location.href = `http://localhost:5173/entrada/${entryId}/${versionId}`; 
    //deberia redireccionar con el versionID de la ultima version, todavia no lo he hecho
  }

  return (
      <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
    
        {/* Formulario para editar Entrada */}
        <Paper elevation={3} sx={{ p: 2, mb: 4 }}>
          <Typography variant="h4" gutterBottom>
            Editar entrada
          </Typography>
          <form id="miFormulario" ref={formRef} onSubmit={subirVersion}>
            <Grid2 container spacing={2}>
              Contenido: 
              <Grid2 item xs={12}>
                <Input
                  id="content"
                  name="content"
                  label="Contenido"
                  value={version.content}
                  onChange={handleChange}
                  multiline
                  fullWidth
                />
              </Grid2>
              Editor:
              <Grid2 item xs={12} sm={6} md={4}>
                <Input
                  id="editor"
                  name="editor"
                  label="Editor"
                  multiline
                  fullWidth
                />
              </Grid2>
              <Grid2 item xs={12} md={4}>
                <Button
                  type="submit"
                  variant="contained"
                  color="primary"
                  fullWidth
                  sx={{ height: "100%" }}
                >
                  Enviar
                </Button>
              </Grid2>
            </Grid2>
          </form>
        </Paper>
      </Container>
  );
}

export default EditarEntradaPage;
