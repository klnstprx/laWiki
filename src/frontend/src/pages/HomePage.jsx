import { useEffect, useState } from "react";
import {
  Container,
  Typography,
  Alert,
  Grid2,
} from "@mui/material";
import { getAllWikis } from "../api/WikiApi.js";
import WikiCard from '../components/WikiCard.jsx';

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
      {error && <Alert severity="error">{error}</Alert>}
      {!error && wikis.length > 0 ? (
        <Grid2
          container
          spacing={4}
          sx={{
            display: "flex",
            justifyContent: "center",
            alignItems: "center",
          }}
        >
          {wikis.map((wiki) => (
            <Grid2 key={wiki.id} xs={12} sm={6} md={4}>
              <WikiCard wiki={wiki} />
            </Grid2>
          ))}
        </Grid2>
      ) : null}
    </Container>
  );
}

export default HomePage;