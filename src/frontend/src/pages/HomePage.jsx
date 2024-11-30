import { useEffect, useState } from "react";
import {
  Container,
  Card,
  CardMedia,
  CardContent,
  Typography,
  Alert,
  Grid2,
  CardActionArea,
} from "@mui/material";
import { Link } from "react-router-dom";
import { getAllWikis } from "../api/WikiApi.js";

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
        <Grid2 container spacing={4}>
          {wikis.map((wiki) => (
            <Grid2 key={wiki.id} xs={12} sm={6} md={4}>
              <Card>
                <CardActionArea component={Link} to={`/wiki/${wiki.id}`}>
                  <CardMedia
                    component="img"
                    height="140"
                    image="https://via.placeholder.com/350x140"
                    alt="Imagen de la Wiki"
                  />
                  <CardContent>
                    <Typography gutterBottom variant="h5" component="div">
                      {wiki.title}
                    </Typography>
                    <Typography variant="body2" color="textSecondary">
                      {wiki.description}
                    </Typography>
                  </CardContent>
                </CardActionArea>
              </Card>
            </Grid2>
          ))}
        </Grid2>
      ) : null}
    </Container>
  );
}

export default HomePage;
