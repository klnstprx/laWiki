import {
  Avatar,
  Typography,
  IconButton,
  Box,
  Card,
  CardContent,
  Stack,
  Rating,
} from "@mui/material";
import DeleteIcon from "@mui/icons-material/Delete";
import PropTypes from "prop-types";
import { getUser } from "../api/AuthApi";
import { useEffect, useState } from "react";

const Comentario = ({ id, content, rating, created_at, author, onDelete }) => {
  const handleDelete = () => {
    onDelete(id);
  };

  const isLoggedIn = !!sessionStorage.getItem('user'); // Verifica si el usuario está logueado
  const [usuario, setUsuario] = useState({}); // add state

  //carga el usuario
    useEffect(() => {
      const fetchUsuario = async () => {
        try {
          const userData = await getUser(author);
          setUsuario(userData);
        } catch (error) {
          console.error("Error fetching user:", error);
        }
      };
  
      if (author) {
        fetchUsuario();
      }
    }, [author]);

  return (
    <Card sx={{ width: "100%" }}>
      <CardContent>
        <Stack direction="row" spacing={2}>
          <Avatar
            src={usuario.picture}
            alt={author}
            sx={{ width: 56, height: 56 }}
          />
          <Box sx={{ flexGrow: 1 }}>
            <Stack
              direction="row"
              justifyContent="space-between"
              alignItems="center"
            >
              <Typography variant="subtitle1" fontWeight="bold">
                {usuario.name}
              </Typography>
              <Typography variant="caption" color="text.secondary">
                {new Date(created_at).toLocaleDateString()}
              </Typography>
            </Stack>
            <Typography
              variant="body1"
              color="text.primary"
              gutterBottom
              sx={{ mt: 1 }}
            >
              {content}
            </Typography>
            <Stack
              direction="row"
              justifyContent="space-between"
              alignItems="center"
            >
              <Rating name="read-only" value={rating} readOnly size="small" />
              {isLoggedIn && (sessionStorage.getItem("role") != "redactor") && ( 
              <IconButton color="error" onClick={handleDelete}>
                <DeleteIcon />
              </IconButton>
              )}
            </Stack>
          </Box>
        </Stack>
      </CardContent>
    </Card>
  );
};

Comentario.propTypes = {
  id: PropTypes.string.isRequired,
  content: PropTypes.string.isRequired,
  rating: PropTypes.number.isRequired,
  created_at: PropTypes.string.isRequired,
  author: PropTypes.string.isRequired,
  onDelete: PropTypes.func.isRequired,
};

export default Comentario;
