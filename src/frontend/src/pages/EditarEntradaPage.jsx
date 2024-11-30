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

function EditarEntradaPage() {
  const { entryId, versionId } = useParams();
  const [version, setVersion] = useState({});
  const [versionError, setVersionError] = useState(null);
  const formRef = useRef(null);
  const navigate = useNavigate();

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
    } else {
      setVersionError("No se proporcionó un ID de versión válido.");
    }
  }, [versionId]);

  // Handler for ReactQuill editor change
  const handleEditorChange = (content, delta, source, editor) => {
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

  // Handler to submit the version
  async function subirVersion(event) {
    event.preventDefault();

    const jsonData = {
      content: version.content,
      editor: version.editor,
      entry_id: entryId,
    };

    console.log("Submitting version:", jsonData); // Debugging

    try {
      const newVersion = await postVersion(jsonData);
      navigate(`/entrada/${entryId}`);
    } catch (error) {
      console.error("Error posting version:", error);
    }
  }

  return (
    <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
      <Paper elevation={3} sx={{ p: { xs: 2, md: 4 }, mb: 4 }}>
        <form id="miFormulario" ref={formRef} onSubmit={subirVersion}>
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
    </Container>
  );
}

export default EditarEntradaPage;
