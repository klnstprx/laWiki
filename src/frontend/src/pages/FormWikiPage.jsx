import { useEffect, useState } from 'react';
import { Container, Typography, TextField, Button, Box, Breadcrumbs } from '@mui/material';
import { useNavigate, useParams, Link } from 'react-router-dom';
import { getWiki, postWiki, putWiki } from '../api/WikiApi';
import ConfirmationModal from '../components/ConfirmationModal.jsx';

function FormWikiPage() {
  const { wikiId } = useParams();
  const [wiki, setWiki] = useState({
    title: "",
    description: "",
    category: "",
  });
  const [error, setError] = useState(null);
  const navigate = useNavigate();

  const [isModalOpen, setIsModalOpen] = useState(false);
  const [formErrors, setFormErrors] = useState({
    title: '',
    description: '',
    category: '',
  });

  useEffect(() => {
    if (wikiId) {
      getWiki(wikiId)
        .then((data) => {
          if (data && Object.keys(data).length > 0) {
            setWiki(data);
          } else {
            setError('Wiki no encontrada.');
          }
        })
        .catch(() => setError('Error al obtener los detalles de la wiki.'));
    }
  }, [wikiId]);

  const handleChange = (event) => {
    const { name, value } = event.target;
    setWiki((prevWiki) => ({
      ...prevWiki,
      [name]: value,
    }));
    // Clear error for the field
    setFormErrors((prevErrors) => ({
      ...prevErrors,
      [name]: '',
    }));
  };

  const validate = () => {
    let isValid = true;
    let errors = { title: '', description: '', category: '' };

    if (!wiki.title.trim()) {
      errors.title = 'El título es obligatorio.';
      isValid = false;
    }
    if (!wiki.description.trim()) {
      errors.description = 'La descripción es obligatoria.';
      isValid = false;
    }
    if (!wiki.category.trim()) {
      errors.category = 'La categoría es obligatoria.';
      isValid = false;
    }

    setFormErrors(errors);
    return isValid;
  };

  const handleSubmit = async () => {
    try {
      if (wikiId) {
        await putWiki(wikiId, wiki);
      } else {
        await postWiki(wiki);
      }
      navigate('/');
    } catch (error) {
      console.error('Error al guardar la wiki:', error);
      setError('Error al guardar la wiki.');
    }
  };

  const onSubmit = (event) => {
    event.preventDefault();
    if (validate()) {
      setIsModalOpen(true);
    }
  };

  return (
    <Container maxWidth="md" sx={{ mt: 4, mb: 4 }}>
      <Breadcrumbs aria-label="breadcrumb" sx={{ mb: 2 }}>
        <Typography color="textPrimary" component={Link} to="/">
          Inicio
        </Typography>
        <Typography color="textPrimary" component={Link} to={`/wiki/${wikiId}`}>
          {wiki.title}
        </Typography>
      </Breadcrumbs>

      <Typography variant="h4" gutterBottom>
        {wikiId ? "Editar Wiki" : "Crear Nueva Wiki"}
      </Typography>
      {error && (
        <Typography variant="body1" color="error" gutterBottom>
          {error}
        </Typography>
      )}
      <Box component="form" onSubmit={onSubmit} noValidate sx={{ mt: 1 }}>
        <TextField
          margin="normal"
          required
          fullWidth
          id="title"
          name="title"
          label="Título"
          value={wiki.title}
          onChange={handleChange}
          variant="outlined"
          error={!!formErrors.title}
          helperText={formErrors.title}
        />
        <TextField
          margin="normal"
          required
          fullWidth
          id="description"
          name="description"
          label="Descripción"
          value={wiki.description}
          onChange={handleChange}
          variant="outlined"
          multiline
          rows={4}
          error={!!formErrors.description}
          helperText={formErrors.description}
        />
        <TextField
          margin="normal"
          required
          fullWidth
          id="category"
          name="category"
          label="Categoría"
          value={wiki.category}
          onChange={handleChange}
          variant="outlined"
          error={!!formErrors.category}
          helperText={formErrors.category}
        />
        <Button
          type="submit"
          fullWidth
          variant="contained"
          sx={{ mt: 3, mb: 2 }}
        >
          {wikiId ? "Guardar Cambios" : "Crear Wiki"}
        </Button>
      </Box>

      <ConfirmationModal
        show={isModalOpen}
        handleClose={() => setIsModalOpen(false)}
        handleConfirm={handleSubmit}
        message={`¿Estás seguro de que quieres ${wikiId ? 'guardar los cambios' : 'crear esta wiki'}?`}
      />
    </Container>
  );
}

export default FormWikiPage;