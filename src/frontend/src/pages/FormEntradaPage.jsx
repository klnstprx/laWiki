import { useState } from 'react';
import { Container, Typography, TextField, Button, Box, IconButton } from '@mui/material';
import { useNavigate, useParams } from 'react-router-dom';
import { postEntry } from '../api/EntryApi';
import DeleteIcon from '@mui/icons-material/Delete';

function FormEntradaPage() {
  const { id: wikiId } = useParams();
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [image, setImage] = useState(null);
  const navigate = useNavigate();

  const handleImageChange = (event) => {
    setImage(event.target.files[0]);
  };

  const handleSubmit = async (event) => {
    event.preventDefault();
    const author = "Guest";
    const formData = new FormData();
    formData.append('title', title);
    formData.append('description', description);
    formData.append('author', author);
    formData.append('wiki_id', wikiId);
    if (image) {
      formData.append('image', image);
    }

    try {
      await postEntry(formData);
      navigate(`/wiki/${wikiId}`);
    } catch (error) {
      console.error('Error al crear la entrada:', error);
    }
  };

  return (
    <Container maxWidth="md">
      <Typography variant="h4" gutterBottom>
        Crear Nueva Entrada
      </Typography>
      <Box component="form" onSubmit={handleSubmit} noValidate sx={{ mt: 1 }}>
        <TextField
          margin="normal"
          required
          fullWidth
          id="title"
          label="Título"
          name="title"
          autoComplete="title"
          autoFocus
          value={title}
          onChange={(e) => setTitle(e.target.value)}
        />
        <TextField
          margin="normal"
          required
          fullWidth
          id="description"
          label="Descripción"
          name="description"
          autoComplete="description"
          value={description}
          onChange={(e) => setDescription(e.target.value)}
        />
        <Button
          variant="contained"
          component="label"
          sx={{ mt: 2 }}
        >
          Añadir Imagen
          <input
            id="image-input"
            type="file"
            hidden
            accept="image/*"
            onChange={handleImageChange}
          />
        </Button>
        {image && (
          <Box sx={{ display: 'flex', alignItems: 'center', mt: 2 }}>
            <Typography variant="body2">
              {image.name}
            </Typography>
            <IconButton onClick={() => setImage(null)} sx={{ ml: 1 }}>
              <DeleteIcon />
            </IconButton>
          </Box>
        )}
        <Button
          type="submit"
          fullWidth
          variant="contained"
          color="primary"
          sx={{ mt: 3, mb: 2 }}
        >
          Crear Entrada
        </Button>
      </Box>
    </Container>
  );
}

export default FormEntradaPage;