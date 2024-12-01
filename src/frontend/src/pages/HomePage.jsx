import { useEffect, useState } from "react";
import { Container, Typography, Alert, Button } from "@mui/material";

import Grid from "@mui/joy/Grid";
import { getAllWikis } from "../api/WikiApi.js";
import WikiCard from "../components/WikiCard.jsx";
import { Link } from "react-router-dom";

function HomePage() {
  const [wikis, setWikis] = useState([]);
  const [error, setError] = useState(null);

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

  return (
    <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
      <Typography variant="h2" gutterBottom>
        Wikis
      </Typography>

      <Button
        variant="contained"
        color="primary"
        component={Link}
        to="/wiki/form"
        sx={{ mb: 2 }}
      >
        Crear Wiki
      </Button>

      {error && <Alert severity="error">{error}</Alert>}
      {!error && wikis.length > 0 ? (
        <Grid
          container
          spacing={4}
          sx={{
            display: "flex",
            justifyContent: "center",
            alignItems: "center",
          }}
        >
          {wikis.map((wiki) => (
            <Grid
              key={wiki.id}
              xs={12}
              sm={6}
              md={4}
              lg={4}
              sx={{ flexBasis: "30%", maxWidth: "30%" }}
            >
              <WikiCard wiki={wiki} />
            </Grid>
          ))}
        </Grid>
      ) : null}
    </Container>
  );
}

export default HomePage;
