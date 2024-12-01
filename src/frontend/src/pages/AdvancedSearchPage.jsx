import { useState, useEffect } from "react";
import {
  Container,
  TextField,
  Button,
  Typography,
  Paper,
  Rating,
  CircularProgress,
} from "@mui/material";
import Grid from "@mui/joy/Grid";
import { DatePicker } from "@mui/x-date-pickers/DatePicker";
import Autocomplete from "@mui/material/Autocomplete";
import SearchResultsList from "../components/SearchResultsList";
import { searchWikis, getAllWikis } from "../api/WikiApi";
import { searchEntries } from "../api/EntryApi";
import { searchComments } from "../api/CommentApi";
import { searchVersions } from "../api/VersionApi";

const AdvancedSearchPage = () => {
  const [params, setParams] = useState({
    // Wikis
    wikiTitle: "",
    wikiDescription: "",
    wikiCategory: "",
    wikiCreatedAtFrom: null,
    wikiCreatedAtTo: null,
    // Entries
    entryTitle: "",
    entryAuthor: "",
    entryCreatedAtFrom: null,
    entryCreatedAtTo: null,
    wikiName: "",
    wikiID: null,
    // Comments
    commentContent: "",
    commentAuthor: "",
    commentCreatedAtFrom: null,
    commentCreatedAtTo: null,
    commentRating: null,
    // Versions
    versionContent: "",
    versionEditor: "",
    versionCreatedAtFrom: null,
    versionCreatedAtTo: null,
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
      if (params.wikiCreatedAtFrom) {
        wikiSearchParams.createdAtFrom = params.wikiCreatedAtFrom.toISOString();
      }
      if (params.wikiCreatedAtTo) {
        wikiSearchParams.createdAtTo = params.wikiCreatedAtTo.toISOString();
      }

      const entriesSearchParams = {};
      if (params.entryTitle.trim() !== "") {
        entriesSearchParams.title = params.entryTitle.trim();
      }
      if (params.entryAuthor.trim() !== "") {
        entriesSearchParams.author = params.entryAuthor.trim();
      }
      if (params.entryCreatedAtFrom) {
        entriesSearchParams.createdAtFrom =
          params.entryCreatedAtFrom.toISOString();
      }
      if (params.entryCreatedAtTo) {
        entriesSearchParams.createdAtTo = params.entryCreatedAtTo.toISOString();
      }
      if (params.wikiID) {
        entriesSearchParams.wikiID = params.wikiID;
      }

      const commentsSearchParams = {};
      if (params.commentContent.trim() !== "") {
        commentsSearchParams.content = params.commentContent.trim();
      }
      if (params.commentAuthor.trim() !== "") {
        commentsSearchParams.author = params.commentAuthor.trim();
      }
      if (params.commentCreatedAtFrom) {
        commentsSearchParams.createdAtFrom =
          params.commentCreatedAtFrom.toISOString();
      }
      if (params.commentCreatedAtTo) {
        commentsSearchParams.createdAtTo =
          params.commentCreatedAtTo.toISOString();
      }
      if (params.commentRating != null) {
        commentsSearchParams.rating = params.commentRating;
      }

      const versionsSearchParams = {};
      if (params.versionContent.trim() !== "") {
        versionsSearchParams.content = params.versionContent.trim();
      }
      if (params.versionEditor.trim() !== "") {
        versionsSearchParams.editor = params.versionEditor.trim();
      }
      if (params.versionCreatedAtFrom) {
        versionsSearchParams.createdAtFrom =
          params.versionCreatedAtFrom.toISOString();
      }
      if (params.versionCreatedAtTo) {
        versionsSearchParams.createdAtTo =
          params.versionCreatedAtTo.toISOString();
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

  const [wikisList, setWikisList] = useState([]);

  useEffect(() => {
    getAllWikis()
      .then((data) => {
        if (data && Array.isArray(data)) {
          setWikisList(data);
        }
      })
      .catch((error) => {
        console.error("Error fetching wikis:", error);
      });
  }, []);

  const isSearchDisabled = Object.values(params).every(
    (value) => value === "" || value === null,
  );

  return (
    <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
      <Typography variant="h4" gutterBottom>
        Búsqueda Avanzada
      </Typography>
      <Paper elevation={3} sx={{ padding: 2, mb: 4 }}>
        <Grid container spacing={2}>
          {/* Wiki Search Fields */}
          <Grid xs={12} sm={6}>
            <Typography variant="h6">Wiki</Typography>

            <TextField
              label="Título de la Wiki"
              name="wikiTitle"
              value={params.wikiTitle}
              onChange={handleInputChange}
              fullWidth
              margin="dense"
            />
            <TextField
              label="Descripción de la Wiki"
              name="wikiDescription"
              value={params.wikiDescription}
              onChange={handleInputChange}
              fullWidth
              margin="dense"
            />
            <TextField
              label="Categoría de la Wiki"
              name="wikiCategory"
              value={params.wikiCategory}
              onChange={handleInputChange}
              fullWidth
              margin="dense"
            />

            <Typography variant="body2" sx={{ marginTop: 2 }}>
              Fecha de Creación de la Wiki
            </Typography>
            <DatePicker
              label="Desde"
              value={params.wikiCreatedAtFrom}
              onChange={(newValue) => {
                setParams((prevParams) => ({
                  ...prevParams,
                  wikiCreatedAtFrom: newValue,
                }));
              }}
              sx={{ mb: 2, mx: 1 }}
            />
            <DatePicker
              label="Hasta"
              value={params.wikiCreatedAtTo}
              onChange={(newValue) => {
                setParams((prevParams) => ({
                  ...prevParams,
                  wikiCreatedAtTo: newValue,
                }));
              }}
              sx={{ mb: 2, mx: 1 }}
            />
          </Grid>

          {/* Entry Search Fields */}
          <Grid xs={12} sm={6}>
            <Typography variant="h6">Entrada</Typography>

            <TextField
              label="Título de la Entrada"
              name="entryTitle"
              value={params.entryTitle}
              onChange={handleInputChange}
              fullWidth
              margin="dense"
            />
            <TextField
              label="Autor de la Entrada"
              name="entryAuthor"
              value={params.entryAuthor}
              onChange={handleInputChange}
              fullWidth
              margin="dense"
            />

            <Typography variant="body2" sx={{ marginTop: 2 }}>
              Fecha de Creación de la Entrada
            </Typography>
            <DatePicker
              label="Desde"
              value={params.entryCreatedAtFrom}
              onChange={(newValue) => {
                setParams((prevParams) => ({
                  ...prevParams,
                  entryCreatedAtFrom: newValue,
                }));
              }}
              sx={{ mb: 2, mx: 1 }}
            />
            <DatePicker
              label="Hasta"
              value={params.entryCreatedAtTo}
              onChange={(newValue) => {
                setParams((prevParams) => ({
                  ...prevParams,
                  entryCreatedAtTo: newValue,
                }));
              }}
              sx={{ mb: 2, mx: 1 }}
            />
            <Autocomplete
              options={wikisList}
              getOptionLabel={(option) => option.title}
              onChange={(event, newValue) => {
                setParams((prevParams) => ({
                  ...prevParams,
                  wikiID: newValue ? newValue.id : null,
                  wikiName: newValue ? newValue.title : "",
                }));
              }}
              renderInput={(params) => (
                <TextField
                  {...params}
                  label="Wiki de la Entrada"
                  margin="dense"
                />
              )}
            />
          </Grid>

          {/* Comment Search Fields */}
          <Grid xs={12} sm={6}>
            <Typography variant="h6">Comentario</Typography>

            <TextField
              label="Contenido del Comentario"
              name="commentContent"
              value={params.commentContent}
              onChange={handleInputChange}
              fullWidth
              margin="dense"
            />
            <TextField
              label="Autor del Comentario"
              name="commentAuthor"
              value={params.commentAuthor}
              onChange={handleInputChange}
              fullWidth
              margin="dense"
            />

            <Typography variant="body2" sx={{ marginTop: 2 }}>
              Fecha de Creación del Comentario
            </Typography>
            <DatePicker
              label="Desde"
              value={params.commentCreatedAtFrom}
              onChange={(newValue) => {
                setParams((prevParams) => ({
                  ...prevParams,
                  commentCreatedAtFrom: newValue,
                }));
              }}
              sx={{ mb: 2, mx: 1 }}
            />
            <DatePicker
              label="Hasta"
              value={params.commentCreatedAtTo}
              onChange={(newValue) => {
                setParams((prevParams) => ({
                  ...prevParams,
                  commentCreatedAtTo: newValue,
                }));
              }}
              sx={{ mb: 2, mx: 1 }}
            />
            <Typography variant="body2">Valoración del Comentario</Typography>
            <Rating
              name="commentRating"
              value={params.commentRating}
              onChange={(event, newValue) => {
                setParams((prevParams) => ({
                  ...prevParams,
                  commentRating: newValue,
                }));
              }}
              size="large"
            />
          </Grid>

          {/* Version Search Fields */}
          <Grid xs={12} sm={6}>
            <Typography variant="h6">Versión</Typography>

            <TextField
              label="Contenido de la Versión"
              name="versionContent"
              value={params.versionContent}
              onChange={handleInputChange}
              fullWidth
              margin="dense"
            />
            <TextField
              label="Editor de la Versión"
              name="versionEditor"
              value={params.versionEditor}
              onChange={handleInputChange}
              fullWidth
              margin="dense"
            />

            <Typography variant="body2" sx={{ marginTop: 2 }}>
              Fecha de Creación de la Versión
            </Typography>
            <DatePicker
              label="Desde"
              value={params.versionCreatedAtFrom}
              onChange={(newValue) => {
                setParams((prevParams) => ({
                  ...prevParams,
                  versionCreatedAtFrom: newValue,
                }));
              }}
              sx={{ mb: 2, mx: 1 }}
            />
            <DatePicker
              label="Hasta"
              value={params.versionCreatedAtTo}
              onChange={(newValue) => {
                setParams((prevParams) => ({
                  ...prevParams,
                  versionCreatedAtTo: newValue,
                }));
              }}
              sx={{ mb: 2, mx: 1 }}
            />
          </Grid>
        </Grid>

        <Button
          variant="contained"
          color="primary"
          onClick={handleSearch}
          sx={{ marginTop: 3 }}
          fullWidth
          disabled={isSearchDisabled}
        >
          Buscar
        </Button>
      </Paper>
      {error && (
        <Typography variant="body2" color="error" sx={{ marginTop: 2 }}>
          {error}
        </Typography>
      )}
      {loading ? (
        <div style={{ display: "flex", alignItems: "center", padding: "10px" }}>
          <CircularProgress size={24} />
          <Typography variant="body2" sx={{ marginLeft: 2 }}>
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
