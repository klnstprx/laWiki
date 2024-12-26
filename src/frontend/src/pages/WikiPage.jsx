import { useEffect, useState } from "react";
import { useParams, Link, useNavigate } from "react-router-dom";
import {
  Container,
  Paper,
  Typography,
  Button,
  Alert,
  Box,
  Breadcrumbs,
  Pagination,
  Menu,
  MenuItem,
} from "@mui/material";

import Grid from "@mui/joy/Grid";
import { deleteEntry, searchEntries } from "../api/EntryApi.js";
import { getWiki, deleteWiki, translateWiki } from "../api/WikiApi.js";
import EntradaCard from "../components/EntradaCard.jsx";
import { useToast } from "../context/ToastContext.jsx";
import ConfirmationModal from "../components/ConfirmationModal.jsx";

const availableLanguages = [
  { code: "AR", name: "Arabic" },
  { code: "BG", name: "Bulgarian" },
  { code: "CS", name: "Czech" },
  { code: "DA", name: "Danish" },
  { code: "DE", name: "German" },
  { code: "EL", name: "Greek" },
  { code: "EN", name: "English" },
  { code: "EN-GB", name: "English (British)" },
  { code: "EN-US", name: "English (American)" },
  { code: "ES", name: "Spanish" },
  { code: "ET", name: "Estonian" },
  { code: "FI", name: "Finnish" },
  { code: "FR", name: "French" },
  { code: "HU", name: "Hungarian" },
  { code: "ID", name: "Indonesian" },
  { code: "IT", name: "Italian" },
  { code: "JA", name: "Japanese" },
  { code: "KO", name: "Korean" },
  { code: "LT", name: "Lithuanian" },
  { code: "LV", name: "Latvian" },
  { code: "NB", name: "Norwegian Bokmål" },
  { code: "NL", name: "Dutch" },
  { code: "PL", name: "Polish" },
  { code: "PT", name: "Portuguese" },
  { code: "PT-BR", name: "Portuguese (Brazilian)" },
  { code: "PT-PT", name: "Portuguese (all Portuguese variants excluding Brazilian Portuguese)" },
  { code: "RO", name: "Romanian" },
  { code: "RU", name: "Russian" },
  { code: "SK", name: "Slovak" },
  { code: "SL", name: "Slovenian" },
  { code: "SV", name: "Swedish" },
  { code: "TR", name: "Turkish" },
  { code: "UK", name: "Ukrainian" },
  { code: "ZH", name: "Chinese" },
  { code: "ZH-HANS", name: "Chinese (simplified)" },
  { code: "ZH-HANT", name: "Chinese (traditional)" },
];

function WikiPage() {
  const [wiki, setWiki] = useState({});
  const [entradas, setEntradas] = useState([]);
  const [error, setError] = useState(null);
  const { showToast } = useToast();
  const navigate = useNavigate();
  const [currentPage, setCurrentPage] = useState(1);
  const itemsPerPage = 6;

  const { id } = useParams();

  const handlePageChange = (event, value) => {
    setCurrentPage(value);
  };

  const startIndex = (currentPage - 1) * itemsPerPage;
  const selectedEntradas = entradas.slice(
    startIndex,
    startIndex + itemsPerPage,
  );

  const [isModalOpen, setIsModalOpen] = useState(false);
  const [selectedOption, setSelectedOption] = useState(null);
  const [anchorEl, setAnchorEl] = useState(null);

  useEffect(() => {
    getWiki(id)
      .then((data) => {
        if (data && Object.keys(data).length > 0) {
          setWiki(data);
        } else {
          setError("Wiki no encontrada.");
        }
      })
      .catch((err) => setError(err.message));
  }, [id]);

  useEffect(() => {
    searchEntries({ wikiID: id })
      .then((data) => {
        if (data && Array.isArray(data)) {
          setEntradas(data);
        } else {
          setEntradas([]);
        }
      })
      .catch((err) => setError(err.message));
  }, [id]);

  const handleDeleteEntry = async (entryID) => {
    try {
      await deleteEntry(entryID);
      setEntradas((prevEntries) =>
        prevEntries.filter((entry) => entry.id !== entryID),
      );
      showToast("Comentario eliminado correctamente", "success");
    } catch (error) {
      console.error("Error al eliminar el comentario:", error);
      showToast("Error al eliminar el comentario", "error");
    }
  };

  const handleDeleteWiki = async () => {
    try {
      await deleteWiki(id);
      showToast("Wiki eliminada correctamente", "success");
      navigate("/");
    } catch (error) {
      console.error("Error al eliminar la wiki:", error);
      showToast("Error al eliminar la wiki", "error");
    }
  };

  const handleDropdownClick = (event) => {
    setAnchorEl(event.currentTarget);
  };

  const handleDropdownClose = () => {
    setAnchorEl(null);
  };

  const handleOptionSelect = (option) => {
    setSelectedOption(option);
    setIsModalOpen(true);
    setAnchorEl(null);
  };

  const handleTranslateWiki = async () => {
    try {
      await translateWiki(id, selectedOption);
      showToast(`Wiki traducida a ${selectedOption} correctamente`, "success");
    } catch (error) {
      console.error("Error al traducir la wiki:", error);
      showToast("Error al traducir la wiki", "error");
    }
    setIsModalOpen(false);
  };


  return (
    <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
      <Breadcrumbs sx={{ mb: 2 }}>
        <Typography
          color="textPrimary"
          component={Link}
          to="/"
          className="breadcrumb-link"
        >
          Inicio
        </Typography>
        <Typography color="textPrimary" className="breadcrumb-active">
          {wiki.title}
        </Typography>
      </Breadcrumbs>

      {error && (
        <Alert severity="error" sx={{ mb: 2 }}>
          {error}
        </Alert>
      )}

      {!error && wiki && Object.keys(wiki).length > 0 && (
        <>
          {/* Page Header */}
          <Paper
            elevation={3}
            sx={{ p: 2,  mb: 4, textAlign: "center", borderRadius: 1 }}
          >
            <Typography variant="h3" component="h1" sx={{ m: 0 }}>
              {wiki.title}
            </Typography>
            <Typography variant="h6" gutterBottom>
              <strong>Descripción:</strong> {wiki.description}
            </Typography>
            <Typography variant="h6" gutterBottom>
              <strong>Categoría:</strong> {wiki.category}
            </Typography>
          </Paper>

          {/* Entradas */}
          <Paper elevation={3} sx={{ p: 3, mb: 4, borderRadius: 1 }}>
            <Typography
              variant="h4"
              component="h2"
              sx={{ borderBottom: "1px solid", pb: 1, mb: 2 }}
            >
              Entradas
            </Typography>
            {selectedEntradas && selectedEntradas.length > 0 ? (
              <Grid container spacing={2}>
                {selectedEntradas.map((entrada) => (
                  <Grid xs={12} sm={6} md={4} key={entrada.id}>
                    <EntradaCard
                      id={entrada.id}
                      title={entrada.title}
                      author={entrada.author}
                      createdAt={entrada.created_at}
                      onDelete={handleDeleteEntry}
                    />
                  </Grid>
                ))}
              </Grid>
            ) : (
              <Typography>No entries available</Typography>
            )}
            <Pagination
              count={Math.ceil(entradas.length / itemsPerPage)}
              page={currentPage}
              onChange={handlePageChange}
              sx={{ mt: 4, display: "flex", justifyContent: "center" }}
            />
          </Paper>

          {/* Buttons */}
          <Box sx={{ display: "flex", justifyContent: "space-between" }}>
            <Button
              component={Link}
              to={`/crear-entrada/${id}`}
              variant="contained"
              color="primary"
              sx={{ mt: 2 }}
            >
              Crear Nueva Entrada
            </Button>

            <Box>
              <Button
                component={Link}
                to={`/wiki/form/${id}`}
                variant="outlined"
                color="primary"
                sx={{ mt: 2, mr: 2 }}
              >
                Editar Wiki
              </Button>
              <Button
                variant="contained"
                color="error"
                sx={{ mt: 2 }}
                onClick={() => setIsModalOpen(true)}
              >
                Borrar Wiki
              </Button>
              <Button
                variant="contained"
                color="primary"
                sx={{ mt: 2, ml: 2 }}
                onClick={handleDropdownClick}
              >
                Traducciones
              </Button>
              <Menu
                anchorEl={anchorEl}
                open={Boolean(anchorEl)}
                onClose={handleDropdownClose}
              >
                {availableLanguages.map((lang) => (
                  <MenuItem key={lang.code} onClick={() => handleOptionSelect(lang.code)}>
                    {lang.name}
                  </MenuItem>
                ))}
              </Menu>
            </Box>
          </Box>

          {/* Confirmation Modal */}
          <ConfirmationModal
            show={isModalOpen}
            handleClose={() => setIsModalOpen(false)}
            handleConfirm={selectedOption ? handleTranslateWiki : handleDeleteWiki}
            message={selectedOption ? 
              `¿Estás seguro de que quieres traducir esta wiki a ${selectedOption}?` 
              : "¿Estás seguro de que quieres borrar esta wiki?"}
          />

        </>
      )}
    </Container>
  );
}

export default WikiPage;
