import { useState } from "react";
import {
  Container,
  TextField,
  Button,
  Typography,
  Paper,
  CircularProgress,
  Grid2,
} from "@mui/material";
import SearchResultsList from "../components/SearchResultsList";
import { searchWikis } from "../api/WikiApi";
import { searchEntries } from "../api/EntryApi";
import { searchComments } from "../api/CommentApi";
import { searchVersions } from "../api/VersionApi";

const AdvancedSearchPage = () => {
  const [params, setParams] = useState({
    wikiTitle: "",
    wikiCategory: "",
    entryTitle: "",
    entryAuthor: "",
    commentContent: "",
    commentAuthor: "",
    versionContent: "",
    versionEditor: "",
  });

  const [searchResults, setSearchResults] = useState({
    wikis: [],
    entries: [],
    comments: [],
    versions: [],
  });

  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setParams((prevParams) => ({
      ...prevParams,
      [name]: value,
    }));
  };

  const handleSearch = async () => {
    setLoading(true);
    setError("");
    try {
      // Build search parameters for each entity
      const wikiSearchParams = {};
      if (params.wikiTitle.trim() !== "") {
        wikiSearchParams.title = params.wikiTitle.trim();
      }
      if (params.wikiCategory.trim() !== "") {
        wikiSearchParams.category = params.wikiCategory.trim();
      }

      const entriesSearchParams = {};
      if (params.entryTitle.trim() !== "") {
        entriesSearchParams.title = params.entryTitle.trim();
      }
      if (params.entryAuthor.trim() !== "") {
        entriesSearchParams.author = params.entryAuthor.trim();
      }

      const commentsSearchParams = {};
      if (params.commentContent.trim() !== "") {
        commentsSearchParams.content = params.commentContent.trim();
      }
      if (params.commentAuthor.trim() !== "") {
        commentsSearchParams.author = params.commentAuthor.trim();
      }

      const versionsSearchParams = {};
      if (params.versionContent.trim() !== "") {
        versionsSearchParams.content = params.versionContent.trim();
      }
      if (params.versionEditor.trim() !== "") {
        versionsSearchParams.editor = params.versionEditor.trim();
      }

      const wikisPromise =
        Object.keys(wikiSearchParams).length > 0
          ? searchWikis(wikiSearchParams)
          : Promise.resolve([]);

      const entriesPromise =
        Object.keys(entriesSearchParams).length > 0
          ? searchEntries(entriesSearchParams)
          : Promise.resolve([]);

      const commentsPromise =
        Object.keys(commentsSearchParams).length > 0
          ? searchComments(commentsSearchParams)
          : Promise.resolve([]);

      const versionsPromise =
        Object.keys(versionsSearchParams).length > 0
          ? searchVersions(versionsSearchParams)
          : Promise.resolve([]);

      const [wikis, entries, comments, versions] = await Promise.all([
        wikisPromise,
        entriesPromise,
        commentsPromise,
        versionsPromise,
      ]);

      setSearchResults({
        wikis: Array.isArray(wikis) ? wikis : [],
        entries: Array.isArray(entries) ? entries : [],
        comments: Array.isArray(comments) ? comments : [],
        versions: Array.isArray(versions) ? versions : [],
      });
    } catch (error) {
      console.error("Error during advanced search:", error);
      setSearchResults({
        wikis: [],
        entries: [],
        comments: [],
        versions: [],
      });
      setError(
        "Ocurrió un error durante la búsqueda. Por favor, inténtelo de nuevo.",
      );
    } finally {
      setLoading(false);
    }
  };

  const isSearchDisabled =
    !params.wikiTitle &&
    !params.wikiCategory &&
    !params.entryTitle &&
    !params.entryAuthor &&
    !params.commentContent &&
    !params.commentAuthor &&
    !params.versionContent &&
    !params.versionEditor;

  return (
    <Container maxWidth="md" sx={{ mt: 4, mb: 4 }}>
      <Typography variant="h4" gutterBottom>
        Búsqueda Avanzada
      </Typography>
      <Paper elevation={3} sx={{ p: 3, mb: 4 }}>
        <Grid2 container spacing={2}>
          {/* Wiki Search Fields */}
          <Grid2 xs={12} sm={6}>
            <Typography variant="h6">Wiki</Typography>
            <TextField
              label="Título de la Wiki"
              name="wikiTitle"
              value={params.wikiTitle}
              onChange={handleInputChange}
              fullWidth
              margin="normal"
            />
            <TextField
              label="Categoría de la Wiki"
              name="wikiCategory"
              value={params.wikiCategory}
              onChange={handleInputChange}
              fullWidth
              margin="normal"
            />
          </Grid2>

          {/* Entry Search Fields */}
          <Grid2 xs={12} sm={6}>
            <Typography variant="h6">Entrada</Typography>
            <TextField
              label="Título de la Entrada"
              name="entryTitle"
              value={params.entryTitle}
              onChange={handleInputChange}
              fullWidth
              margin="normal"
            />
            <TextField
              label="Autor de la Entrada"
              name="entryAuthor"
              value={params.entryAuthor}
              onChange={handleInputChange}
              fullWidth
              margin="normal"
            />
          </Grid2>

          {/* Comment Search Fields */}
          <Grid2 xs={12} sm={6}>
            <Typography variant="h6">Comentario</Typography>
            <TextField
              label="Contenido del Comentario"
              name="commentContent"
              value={params.commentContent}
              onChange={handleInputChange}
              fullWidth
              margin="normal"
            />
            <TextField
              label="Autor del Comentario"
              name="commentAuthor"
              value={params.commentAuthor}
              onChange={handleInputChange}
              fullWidth
              margin="normal"
            />
          </Grid2>

          {/* Version Search Fields */}
          <Grid2 xs={12} sm={6}>
            <Typography variant="h6">Versión</Typography>
            <TextField
              label="Contenido de la Versión"
              name="versionContent"
              value={params.versionContent}
              onChange={handleInputChange}
              fullWidth
              margin="normal"
            />
            <TextField
              label="Editor de la Versión"
              name="versionEditor"
              value={params.versionEditor}
              onChange={handleInputChange}
              fullWidth
              margin="normal"
            />
          </Grid2>
        </Grid2>

        <Button
          variant="contained"
          color="primary"
          onClick={handleSearch}
          sx={{ mt: 2 }}
          fullWidth
          disabled={isSearchDisabled}
        >
          Buscar
        </Button>
      </Paper>
      {error && (
        <Typography variant="body2" color="error" sx={{ mt: 2 }}>
          {error}
        </Typography>
      )}
      {loading ? (
        <div style={{ display: "flex", alignItems: "center", padding: "10px" }}>
          <CircularProgress size={24} />
          <Typography variant="body2" sx={{ ml: 2 }}>
            Cargando...
          </Typography>
        </div>
      ) : (
        <SearchResultsList results={searchResults} />
      )}
    </Container>
  );
};

export default AdvancedSearchPage;
