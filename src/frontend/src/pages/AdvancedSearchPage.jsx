import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { Container, TextField, Button, Typography, Grid2 } from "@mui/material";

const AdvancedSearchPage = () => {
  const navigate = useNavigate();
  const [params, setParams] = useState({
    wikiTitle: "",
    wikiCategory: "",
    entryTitle: "",
    entryAuthor: "",
    commentContent: "",
    commentAuthor: "",
    versionContent: "",
    versionEditor: "",
  });

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setParams((prevParams) => ({
      ...prevParams,
      [name]: value,
    }));
  };

  const handleSearch = () => {
    const urlParams = new URLSearchParams();

    // Add wiki search parameters
    if (params.wikiTitle) urlParams.append("wikiTitle", params.wikiTitle);
    if (params.wikiCategory)
      urlParams.append("wikiCategory", params.wikiCategory);

    // Add entry search parameters
    if (params.entryTitle) urlParams.append("entryTitle", params.entryTitle);
    if (params.entryAuthor) urlParams.append("entryAuthor", params.entryAuthor);

    // Add comment search parameters
    if (params.commentContent)
      urlParams.append("commentContent", params.commentContent);
    if (params.commentAuthor)
      urlParams.append("commentAuthor", params.commentAuthor);

    // Add version search parameters
    if (params.versionContent)
      urlParams.append("versionContent", params.versionContent);
    if (params.versionEditor)
      urlParams.append("versionEditor", params.versionEditor);

    navigate(`/search?${urlParams.toString()}`);
  };

  return (
    <Container maxWidth="md" sx={{ mt: 4 }}>
      <Typography variant="h4" gutterBottom>
        Advanced Search
      </Typography>
      <form>
        <Grid2 container spacing={3}>
          {/* Wiki Search Fields */}
          <Grid2 item xs={12} sm={6}>
            <Typography variant="h6">Wiki </Typography>
            <Grid2 container spacing={2}>
              <Grid2 item xs={12} sm={6}>
                <TextField
                  label="Wiki Title"
                  name="wikiTitle"
                  value={params.wikiTitle}
                  onChange={handleInputChange}
                  fullWidth
                />
              </Grid2>
              <Grid2 item xs={12} sm={6}>
                <TextField
                  label="Wiki Category"
                  name="wikiCategory"
                  value={params.wikiCategory}
                  onChange={handleInputChange}
                  fullWidth
                />
              </Grid2>
            </Grid2>
          </Grid2>

          {/* Entry Search Fields */}
          <Grid2 item xs={12} sm={6}>
            <Typography variant="h6">Entrada </Typography>
            <Grid2 container spacing={2}>
              <Grid2 item xs={12} sm={6}>
                <TextField
                  label="Entry Title"
                  name="entryTitle"
                  value={params.entryTitle}
                  onChange={handleInputChange}
                  fullWidth
                />
              </Grid2>
              <Grid2 item xs={12} sm={6}>
                <TextField
                  label="Entry Author"
                  name="entryAuthor"
                  value={params.entryAuthor}
                  onChange={handleInputChange}
                  fullWidth
                />
              </Grid2>
            </Grid2>
          </Grid2>

          {/* Comment Search Fields */}
          <Grid2 item xs={12} sm={6}>
            <Typography variant="h6">Commentario</Typography>
            <Grid2 container spacing={2}>
              <Grid2 item xs={12} sm={6}>
                <TextField
                  label="Comment Content"
                  name="commentContent"
                  value={params.commentContent}
                  onChange={handleInputChange}
                  fullWidth
                />
              </Grid2>
              <Grid2 item xs={12} sm={6}>
                <TextField
                  label="Comment Author"
                  name="commentAuthor"
                  value={params.commentAuthor}
                  onChange={handleInputChange}
                  fullWidth
                />
              </Grid2>
            </Grid2>
          </Grid2>

          {/* Version Search Fields */}
          <Grid2 item xs={12} sm={6}>
            <Typography variant="h6">Version</Typography>
            <Grid2 container spacing={2}>
              <Grid2 item xs={12} sm={6}>
                <TextField
                  label="Version Content"
                  name="versionContent"
                  value={params.versionContent}
                  onChange={handleInputChange}
                  fullWidth
                />
              </Grid2>
              <Grid2 item xs={12} sm={6}>
                <TextField
                  label="Version Editor"
                  name="versionEditor"
                  value={params.versionEditor}
                  onChange={handleInputChange}
                  fullWidth
                />
              </Grid2>
            </Grid2>
          </Grid2>

          <Button
            variant="contained"
            color="primary"
            onClick={handleSearch}
            fullWidth
          >
            Buscar
          </Button>
        </Grid2>
      </form>
    </Container>
  );
};

export default AdvancedSearchPage;
