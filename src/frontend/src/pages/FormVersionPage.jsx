import { useEffect, useState, useRef } from "react";
import { postVersion, getVersion } from "../api/VersionApi.js";
import { useNavigate, useParams } from "react-router-dom";
import {
  Container,
  Paper,
  Typography,
  Button,
  TextField,
  Box,
} from "@mui/material";
import Grid from "@mui/joy/Grid";
import ReactQuill from "react-quill";
import "react-quill/dist/quill.snow.css";
import ConfirmationModal from "../components/ConfirmationModal.jsx"; // add import

function FormVersionPage() {
  const { entryId, versionId } = useParams();
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
      <Paper elevation={3} sx={{ p: { xs: 2, md: 4 }, mb: 4 }}>
        <form id="miFormulario" ref={formRef} onSubmit={onSubmit}>
          {/* Title and Editor input */}
          <Grid container spacing={2} alignItems="center">
            <Grid xs={12} sm={8}>
              <Typography variant="h4">Editar entrada</Typography>
            </Grid>
            <Grid xs={12} sm={4}>
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
            </Grid>
            <Grid xs={12} sm={4}>
              <TextField
                id="address"
                name="address"
                label="Ubicación (opcional)"
                value={version.address || ""}
                onChange={handleChange}
                variant="outlined"
                fullWidth
              />
            </Grid>
          </Grid>
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
        message={`¿Estás seguro de que deseas ${
          versionId ? "guardar los cambios" : "crear esta versión"
        }?`}
      />
    </Container>
  );
}

export default FormVersionPage;
