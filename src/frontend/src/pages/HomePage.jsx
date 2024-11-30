import { useEffect, useState } from "react";
import {
  Container,
  Grid2,
  Card,
  CardMedia,
  CardContent,
  Typography,
  Alert,
} from "@mui/material";
import { getAllWikis } from "../api/WikiApi.js";
import MainLayout from "../layout/MainLayout.jsx";

function HomePage() {
  const [wikis, setWikis] = useState([]);
  const [error, setError] = useState(null);

  useEffect(() => {
    getAllWikis()
      .then(setWikis)
      .catch((err) => setError(err.message));
  }, []);

  return (
    <MainLayout>
      <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
        <Typography variant="h2" gutterBottom>
          Wikis
        </Typography>
        {error && <Alert severity="error">{error}</Alert>}
        <Grid2 container spacing={4}>
          {wikis.map((wiki) => (
            <Grid2 item key={wiki.id} xs={12} sm={6} md={4}>
              <Card>
                <CardMedia
                  component="img"
                  height="140"
                  image="https://via.placeholder.com/350x140"
                  alt="Wiki Image"
                />
                <CardContent>
                  <Typography gutterBottom variant="h5" component="div">
                    {wiki.title}
                  </Typography>
                  <Typography variant="body2" color="textSecondary">
                    {wiki.description}
                  </Typography>
                </CardContent>
              </Card>
            </Grid2>
          ))}
        </Grid2>
      </Container>
    </MainLayout>
  );
}

export default HomePage;
