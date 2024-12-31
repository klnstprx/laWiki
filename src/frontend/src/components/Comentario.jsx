import {
  Avatar,
  Box,
  Card,
  CardContent,
  IconButton,
  Rating,
  Stack,
  Typography,
} from "@mui/material";
import DeleteIcon from "@mui/icons-material/Delete";
import PropTypes from "prop-types";
import { getUser } from "../api/AuthApi";
import { useEffect, useState } from "react";

const Comentario = (
  { id, content, rating, created_at, authorId, onDelete },
) => {
  const handleDelete = () => {
    onDelete(id);
  };

  // Retrieve the JSON string for 'appUser' from sessionStorage
  const appUserJson = sessionStorage.getItem("appUser");

  // Check if 'appUser' exists in sessionStorage
  const isLoggedIn = !!appUserJson;

  let role = null;

  if (isLoggedIn) {
    const appUser = JSON.parse(appUserJson);

    role = appUser.role;
  }
  const [author, setAuthor] = useState({}); // add state

  //carga el usuario
  useEffect(() => {
    const fetchUsuario = async () => {
      try {
        const userData = await getUser(authorId);
        setAuthor(userData);
      } catch (error) {
        console.error("Error fetching user:", error);
      }
    };

    if (authorId) {
      fetchUsuario();
    }
  }, [authorId]);

  return (
    <Card
      sx={{
        width: "100%",
        mb: 2,
        "&:hover": {
          boxShadow: 6,
        },
        transition: "box-shadow 0.3s",
      }}
    >
      <CardContent>
        <Stack direction="row" spacing={2}>
          <Avatar
            src={author.picture}
            alt={author.name}
            sx={{ width: 56, height: 56 }}
          />
          <Box sx={{ flexGrow: 1 }}>
            <Stack
              direction="row"
              justifyContent="space-between"
              alignItems="center"
            >
              <Typography variant="subtitle1" fontWeight="bold">
                {author.name}
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
              {isLoggedIn && (role != "redactor") && (
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
  authorId: PropTypes.string.isRequired,
  onDelete: PropTypes.func.isRequired,
};

export default Comentario;
