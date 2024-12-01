import { useEffect, useState, useRef } from "react";
import { postVersion, getVersion } from "../api/VersionApi.js";
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
import ReactQuill from "react-quill";
import "react-quill/dist/quill.snow.css";
import ConfirmationModal from '../components/ConfirmationModal.jsx';
import { Link } from "react-router-dom";
import { Breadcrumbs } from "@mui/material";
import { getEntry } from "../api/EntryApi.js";
import { getWiki } from "../api/WikiApi.js";

function FormVersionPage() {
  const { entryId, versionId } = useParams();
  const [entry, setEntry] = useState({});
  const [wiki, setWiki] = useState({});
  const [version, setVersion] = useState({});
  const [versionError, setVersionError] = useState(null);
  const formRef = useRef(null);
  const navigate = useNavigate();
  const [isModalOpen, setIsModalOpen] = useState(false); // add state

  // Fetch the version details
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
          setVersionError("Se produjo un error al obtener la entrada."),
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
          setVersionError("Se produjo un error al obtener la wiki asociada."),
        );
    }
  }, [entry, entryId]);

  // Handler for ReactQuill editor change
  const handleEditorChange = (content) => {
    setVersion((prevVersion) => ({
      ...prevVersion,
      content: content,
    }));
  };

  // Handler for TextField change
  const handleChange = (event) => {
    const { name, value } = event.target;
    setVersion((prevVersion) => ({
      ...prevVersion,
      [name]: value,
    }));
  };

  // Validation function
  const validate = () => {
    if (!version.editor) {
      setVersionError("El campo Editor es obligatorio.");
      return false;
    }
    setVersionError(null);
    return true;
  };

  // Handler to submit the version
  async function handleSubmit(event) {
    event.preventDefault();

    const jsonData = {
      content: version.content,
      editor: version.editor,
      entry_id: entryId,
      address: version.address,
    };

    console.log("Submitting version:", jsonData); // Debugging

    try {
      await postVersion(jsonData);
      navigate(`/entrada/${entryId}`);
    } catch (error) {
      console.error("Error posting version:", error);
    }
  }

  const onSubmit = (event) => {
    event.preventDefault();
    if (validate()) {
      setIsModalOpen(true);
    }
  };

  return (
    <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
      <Breadcrumbs sx={{ mb: 2 }}>
        <Typography color="textPrimary" component={Link} to="/">
          Inicio
        </Typography>
        <Typography color="textPrimary" component={Link} to={`/wiki/${wiki.id}`}>
          {wiki.title}
        </Typography>
        <Typography color="textPrimary" component={Link} to={`/entrada/${entry.id}`}>
          {entry.title}
        </Typography>
      </Breadcrumbs>

      <Paper elevation={3} sx={{ p: { xs: 2, md: 4 }, mb: 4 }}>
        <form id="miFormulario" ref={formRef} onSubmit={onSubmit}>
          {/* Title and Editor input */}
          <Grid2 container spacing={2} alignItems="center">
            <Grid2 item xs={12} sm={8}>
              <Typography variant="h4">Editar entrada</Typography>
            </Grid2>
            <Grid2 item xs={12} sm={4}>
              <TextField
                id="editor"
                name="editor"
                label="Editor"
                value={version.editor || ""}
                onChange={handleChange}
                variant="outlined"
                fullWidth
                required
              />
            </Grid2>
            <Grid2 item xs={12} sm={4}>
              <TextField
                id="address"
                name="address"
                label="Ubicación (opcional)"
                value={version.address || ""}
                onChange={handleChange}
                variant="outlined"
                fullWidth
              />
            </Grid2>
          </Grid2>
          {versionError && (
            <Typography variant="body1" color="error" gutterBottom>
              {versionError}
            </Typography>
          )}
          <Typography variant="subtitle1" sx={{ mt: 3, mb: 1 }}>
            Contenido:
          </Typography>
          <Box sx={{ height: "400px", mb: 2 }}>
            <ReactQuill
              theme="snow"
              value={version.content || ""}
              onChange={handleEditorChange}
              style={{ height: "100%" }}
            />
          </Box>
          <Box sx={{ mt: 5, pt: 1 }} display="flex" justifyContent="flex-end">
            <Button type="submit" variant="contained" color="primary">
              Enviar
            </Button>
          </Box>
        </form>
      </Paper>
      <ConfirmationModal
        show={isModalOpen}
        handleClose={() => setIsModalOpen(false)}
        handleConfirm={handleSubmit}
        message={`¿Estás seguro de que deseas ${versionId ? 'guardar los cambios' : 'crear esta versión'}?`}
      />
    </Container>
  );
}

export default FormVersionPage;
