import { useState, useEffect } from "react";
import List from "@mui/material/List";
import ListItem from "@mui/material/ListItem";
import Alert from "@mui/material/Alert";
import Typography from "@mui/material/Typography";
import { getAllWikis } from "../api";

function Home() {
  const [wikis, setWikis] = useState([]);
  const [error, setError] = useState(null);

  useEffect(() => {
    getAllWikis()
      .then(setWikis)
      .catch((err) => setError(err.message));
  }, []);

  return (
    <div>
      <Typography variant="h1" gutterBottom>
        Wikis
      </Typography>
      {error && <Alert severity="error">{error}</Alert>}
      {wikis.length > 0 ? (
        <List>
          {wikis.map((wiki) => (
            <ListItem key={wiki.id}>
              <Typography variant="h2" gutterBottom>
                {wiki.title}
              </Typography>
              <Typography variant="body1" gutterBottom>
                {wiki.description}
              </Typography>
            </ListItem>
          ))}
        </List>
      ) : (
        <Alert>No wikis found.</Alert>
      )}
    </div>
  );
}

export default Home;
