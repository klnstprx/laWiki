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
} from "@mui/material";
import ReactQuill from "react-quill";
import "react-quill/dist/quill.snow.css";

function EditarEntradaPage() {
  const [version, setVersion] = useState({});
  const [versionError, setVersionError] = useState("");
  const navigate = useNavigate();

  const formRef = useRef(null);
  const { entryId, versionId } = useParams();

  // Load the current version
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
      {/* Form to edit Entry */}
      <Paper elevation={3} sx={{ p: 2, mb: 4 }}>
        <Typography variant="h4" gutterBottom>
          Editar entrada
        </Typography>
        {versionError && (
          <Typography variant="body1" color="error" gutterBottom>
            {versionError}
          </Typography>
        )}
        <form id="miFormulario" ref={formRef} onSubmit={subirVersion}>
          <Grid2 container spacing={2}>
            <Grid2 item xs={12}>
              <Typography variant="subtitle1">Contenido:</Typography>
              <ReactQuill
                theme="snow"
                value={version.content || ""}
                onChange={handleEditorChange}
                style={{ height: "300px", marginBottom: "50px" }}
              />
            </Grid2>
            <Grid2 item xs={12} sm={6} md={4}>
              <TextField
                id="editor"
                name="editor"
                label="Editor"
                value={version.editor || ""}
                onChange={handleChange}
                fullWidth
                variant="outlined"
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
