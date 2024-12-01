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

const AdvancedSearchPage = () => {
  const [params, setParams] = useState({
    wikiTitle: "",
    wikiCategory: "",
    entryTitle: "",
    entryAuthor: "",
  });

  const [searchResults, setSearchResults] = useState({
    wikis: [],
    entries: [],
  });

  const [loading, setLoading] = useState(false);

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setParams((prevParams) => ({
      ...prevParams,
      [name]: value,
    }));
  };
  const [error, setError] = useState("");
  const handleSearch = async () => {
    setLoading(true);
    setError("");
    try {
      // Build wiki search parameters
      const wikiSearchParams = {};
      if (params.wikiTitle.trim() !== "") {
        wikiSearchParams.title = params.wikiTitle.trim();
      }
      if (params.wikiCategory.trim() !== "") {
        wikiSearchParams.category = params.wikiCategory.trim();
      }

      // Build entry search parameters
      const entriesSearchParams = {};
      if (params.entryTitle.trim() !== "") {
        entriesSearchParams.title = params.entryTitle.trim();
      }
      if (params.entryAuthor.trim() !== "") {
        entriesSearchParams.author = params.entryAuthor.trim();
      }

      // Only call the API if there are parameters to search
      const wikisPromise =
        Object.keys(wikiSearchParams).length > 0
          ? searchWikis(wikiSearchParams)
          : Promise.resolve([]);

      const entriesPromise =
        Object.keys(entriesSearchParams).length > 0
          ? searchEntries(entriesSearchParams)
          : Promise.resolve([]);

      const [wikis, entries] = await Promise.all([
        wikisPromise,
        entriesPromise,
      ]);

      console.log("Received wikis:", wikis);
      console.log("Received entries:", entries);

      setSearchResults({
        wikis: Array.isArray(wikis) ? wikis : [],
        entries: Array.isArray(entries) ? entries : [],
      });
    } catch (error) {
      console.error("Error during advanced search:", error);
      setSearchResults({ wikis: [], entries: [] });
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
    !params.entryAuthor;

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
