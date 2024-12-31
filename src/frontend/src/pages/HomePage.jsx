import { useEffect, useState } from "react";
import {
  Alert,
  Breadcrumbs,
  Button,
  Container,
  Pagination,
  Paper,
  Typography,
} from "@mui/material";

import Grid from "@mui/joy/Grid";
import { getAllWikis } from "../api/WikiApi.js";
import WikiCard from "../components/WikiCard.jsx";
import { Link } from "react-router-dom";

function HomePage() {
  const [wikis, setWikis] = useState([]);
  const [error, setError] = useState(null);
  const [currentPage, setCurrentPage] = useState(1);
  const itemsPerPage = 6;

  useEffect(() => {
    getAllWikis()
      .then((data) => {
        if (data && data.length > 0) {
          setWikis(data);
        } else {
          setError("No se encontraron wikis.");
        }
      })
      .catch(() => setError("Se produjo un error al obtener las wikis."));
  }, []);

  const handlePageChange = (event, value) => {
    setCurrentPage(value);
  };

  const startIndex = (currentPage - 1) * itemsPerPage;
  const selectedWikis = wikis.slice(startIndex, startIndex + itemsPerPage);
  // Retrieve the JSON string for 'appUser' from sessionStorage
  const appUserJson = sessionStorage.getItem("appUser");

  // Check if 'appUser' exists in sessionStorage
  const isLoggedIn = !!appUserJson;

  let role = null;

  if (isLoggedIn) {
    const appUser = JSON.parse(appUserJson);

    role = appUser.role;
  }

  return (
    <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
      <Breadcrumbs sx={{ mb: 2 }}>
        <Typography className="breadcrumb-active">Inicio</Typography>
      </Breadcrumbs>

      <Paper sx={{ p: 2, mb: 4, textAlign: "center", borderRadius: 1 }}>
        <Typography variant="h2" gutterBottom>
          Wikis
        </Typography>
        {isLoggedIn && (role != "redactor") && (
          <Button
            variant="contained"
            color="primary"
            component={Link}
            to="/wiki/form"
            sx={{ mb: 2 }}
          >
            Crear Wiki
          </Button>
        )}
      </Paper>

      <Paper sx={{ p: 2, mb: 4, textAlign: "center", borderRadius: 1 }}>
        {error && <Alert severity="error">{error}</Alert>}
        {!error && wikis.length > 0
          ? (
            <>
              <Grid
                container
                spacing={4}
                sx={{
                  display: "flex",
                  justifyContent: "center",
                  alignItems: "top",
                }}
              >
                {selectedWikis.map((wiki) => (
                  <Grid
                    key={wiki.id}
                    xs={12}
                    sm={6}
                    md={4}
                    lg={3}
                  >
                    <WikiCard wiki={wiki} />
                  </Grid>
                ))}
              </Grid>
              <Pagination
                count={Math.ceil(wikis.length / itemsPerPage)}
                page={currentPage}
                onChange={handlePageChange}
                sx={{ mt: 4, display: "flex", justifyContent: "center" }}
              />
            </>
          )
          : null}
      </Paper>
    </Container>
  );
}

export default HomePage;
