import { useEffect, useState } from "react";
import {
  Button,
  CircularProgress,
  Container,
  Paper,
  Rating,
  TextField,
  Typography,
} from "@mui/material";
import Grid from "@mui/joy/Grid";
import { DatePicker } from "@mui/x-date-pickers/DatePicker";
import Autocomplete from "@mui/material/Autocomplete";
import Breadcrumbs from "@mui/material/Breadcrumbs";
import { Link } from "react-router-dom";
import SearchResultsList from "../components/SearchResultsList";
import { getAllWikis, searchWikis } from "../api/WikiApi";
import { searchEntries } from "../api/EntryApi";
import { searchComments } from "../api/CommentApi";
import { searchVersions } from "../api/VersionApi";
import { getUsersByIds, getUsersByName } from "../api/AuthApi";

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
  const [hasSearched, setHasSearched] = useState(false);

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
    setHasSearched(true);

    try {
      // Initialize search parameters
      const wikiSearchParams = {};
      const entriesSearchParams = {};
      const commentsSearchParams = {};
      const versionsSearchParams = {};

      // Map to hold user IDs and names for display purposes
      const userIdToName = {};

      // Fetch user IDs based on input names
      let entryAuthorIds = [];
      if (params.entryAuthor.trim() !== "") {
        const entryAuthors = await getUsersByName(params.entryAuthor.trim());
        if (!entryAuthors || entryAuthors.length === 0) {
          // No matching authors found; set empty results
          setSearchResults({
            wikis: [],
            entries: [],
            comments: [],
            versions: [],
          });
          setLoading(false);
          return;
        }
        entryAuthorIds = entryAuthors.map((user) => {
          userIdToName[user.id] = user.name;
          return user.id;
        });
      }

      // Repeat for comment authors and version editors
      let commentAuthorIds = [];
      if (params.commentAuthor.trim() !== "") {
        const commentAuthors = await getUsersByName(
          params.commentAuthor.trim(),
        );
        if (!commentAuthors || commentAuthors.length === 0) {
          setSearchResults({
            wikis: [],
            entries: [],
            comments: [],
            versions: [],
          });
          setLoading(false);
          return;
        }
        commentAuthorIds = commentAuthors.map((user) => {
          userIdToName[user.id] = user.name;
          return user.id;
        });
      }

      let versionEditorIds = [];
      if (params.versionEditor.trim() !== "") {
        const versionEditors = await getUsersByName(
          params.versionEditor.trim(),
        );
        if (!versionEditors || versionEditors.length === 0) {
          setSearchResults({
            wikis: [],
            entries: [],
            comments: [],
            versions: [],
          });
          setLoading(false);
          return;
        }
        versionEditorIds = versionEditors.map((user) => {
          userIdToName[user.id] = user.name;
          return user.id;
        });
      }

      // Build search parameters for each entity using IDs
      // Wikis
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

      // Entries
      if (params.entryTitle.trim() !== "") {
        entriesSearchParams.title = params.entryTitle.trim();
      }
      if (entryAuthorIds.length > 0) {
        entriesSearchParams.author = entryAuthorIds; // Note that this is an array
      }
      if (params.entryCreatedAtFrom) {
        entriesSearchParams.createdAtFrom = params.entryCreatedAtFrom
          .toISOString();
      }
      if (params.entryCreatedAtTo) {
        entriesSearchParams.createdAtTo = params.entryCreatedAtTo.toISOString();
      }
      if (params.wikiID) {
        entriesSearchParams.wikiID = params.wikiID;
      }

      // Comments
      if (params.commentContent.trim() !== "") {
        commentsSearchParams.content = params.commentContent.trim();
      }
      if (commentAuthorIds.length > 0) {
        commentsSearchParams.author = commentAuthorIds; // This is an array
      }
      if (params.commentCreatedAtFrom) {
        commentsSearchParams.createdAtFrom = params.commentCreatedAtFrom
          .toISOString();
      }
      if (params.commentCreatedAtTo) {
        commentsSearchParams.createdAtTo = params.commentCreatedAtTo
          .toISOString();
      }
      if (params.commentRating != null) {
        commentsSearchParams.rating = params.commentRating;
      }

      // Versions
      if (params.versionContent.trim() !== "") {
        versionsSearchParams.content = params.versionContent.trim();
      }
      if (versionEditorIds.length > 0) {
        versionsSearchParams.editor = versionEditorIds; // This is an array
      }
      if (params.versionCreatedAtFrom) {
        versionsSearchParams.createdAtFrom = params.versionCreatedAtFrom
          .toISOString();
      }
      if (params.versionCreatedAtTo) {
        versionsSearchParams.createdAtTo = params.versionCreatedAtTo
          .toISOString();
      }

      // Perform searches
      const [wikis, entries, comments, versions] = await Promise.all([
        Object.keys(wikiSearchParams).length > 0
          ? searchWikis(wikiSearchParams)
          : Promise.resolve([]),
        Object.keys(entriesSearchParams).length > 0
          ? searchEntries(entriesSearchParams)
          : Promise.resolve([]),
        Object.keys(commentsSearchParams).length > 0
          ? searchComments(commentsSearchParams)
          : Promise.resolve([]),
        Object.keys(versionsSearchParams).length > 0
          ? searchVersions(versionsSearchParams)
          : Promise.resolve([]),
      ]);

      // Collect user IDs from search results
      const allUserIds = new Set();
      entries.forEach((entry) => {
        if (entry.author) allUserIds.add(entry.author);
      });
      comments.forEach((comment) => {
        if (comment.author) allUserIds.add(comment.author);
      });
      versions.forEach((version) => {
        if (version.editor) allUserIds.add(version.editor);
      });

      // Fetch user data for IDs not already in the userIdToName map
      const missingUserIds = Array.from(allUserIds).filter(
        (id) => !userIdToName[id],
      );
      if (missingUserIds.length > 0) {
        const additionalUsers = await getUsersByIds(missingUserIds);
        additionalUsers.forEach((user) => {
          userIdToName[user.id] = user.name;
        });
      }

      // Prepare the results with user names
      const processedEntries = entries.map((entry) => ({
        ...entry,
        authorName: userIdToName[entry.author] || "Usuario desconocido",
      }));

      const processedComments = comments.map((comment) => ({
        ...comment,
        authorName: userIdToName[comment.author] || "Usuario desconocido",
      }));

      const processedVersions = versions.map((version) => ({
        ...version,
        editorName: userIdToName[version.editor] || "Usuario desconocido",
      }));

      // Set the search results
      setSearchResults({
        wikis: Array.isArray(wikis) ? wikis : [],
        entries: Array.isArray(processedEntries) ? processedEntries : [],
        comments: Array.isArray(processedComments) ? processedComments : [],
        versions: Array.isArray(processedVersions) ? processedVersions : [],
      });
    } catch (error) {
      // Handle errors
      console.error("Error during advanced search:", error);
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
      <Breadcrumbs aria-label="breadcrumb" sx={{ mb: 2 }}>
        <Typography component={Link} to="/" className="breadcrumb-link">
          Inicio
        </Typography>
      </Breadcrumbs>
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
              slotProps={{
                inputLabel: {
                  shrink: true,
                },
              }}
            />
            <TextField
              label="Descripción de la Wiki"
              name="wikiDescription"
              value={params.wikiDescription}
              onChange={handleInputChange}
              fullWidth
              margin="dense"
              slotProps={{
                inputLabel: {
                  shrink: true,
                },
              }}
            />
            <TextField
              label="Categoría de la Wiki"
              name="wikiCategory"
              value={params.wikiCategory}
              onChange={handleInputChange}
              fullWidth
              margin="dense"
              slotProps={{
                inputLabel: {
                  shrink: true,
                },
              }}
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
              slotProps={{
                inputLabel: {
                  shrink: true,
                },
              }}
            />
            <TextField
              label="Autor de la Entrada"
              name="entryAuthor"
              value={params.entryAuthor}
              onChange={handleInputChange}
              fullWidth
              margin="dense"
              slotProps={{
                inputLabel: {
                  shrink: true,
                },
              }}
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
                  slotProps={{
                    inputLabel: {
                      shrink: true,
                    },
                  }}
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
              slotProps={{
                inputLabel: {
                  shrink: true,
                },
              }}
            />
            <TextField
              label="Autor del Comentario"
              name="commentAuthor"
              value={params.commentAuthor}
              onChange={handleInputChange}
              fullWidth
              margin="dense"
              slotProps={{
                inputLabel: {
                  shrink: true,
                },
              }}
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
              slotProps={{
                inputLabel: {
                  shrink: true,
                },
              }}
            />
            <TextField
              label="Editor de la Versión"
              name="versionEditor"
              value={params.versionEditor}
              onChange={handleInputChange}
              fullWidth
              margin="dense"
              slotProps={{
                inputLabel: {
                  shrink: true,
                },
              }}
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
              slotProps={{
                inputLabel: {
                  shrink: true,
                },
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
      {loading
        ? (
          <div
            style={{ display: "flex", alignItems: "center", padding: "10px" }}
          >
            <CircularProgress size={24} />
            <Typography variant="body2" sx={{ marginLeft: 2 }}>
              Cargando...
            </Typography>
          </div>
        )
        : (
          hasSearched && <SearchResultsList results={searchResults} />
        )}
    </Container>
  );
};

export default AdvancedSearchPage;
