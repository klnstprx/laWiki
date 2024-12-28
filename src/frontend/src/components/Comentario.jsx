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
import { useAuth } from "../context/AuthContext";
import { Link } from "react-router-dom";

const Comentario = ({ id, content, rating, created_at, author, onDelete }) => {
  const handleDelete = () => {
    onDelete(id);
  };

  const { user } = useAuth();
  const isLoggedIn = !!user;
  const userRole = user?.role || "";

  return (
    <Card sx={{ width: "100%" }}>
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
              <Typography
                component={Link}
                to={`/perfil/${author.id}`}
                variant="subtitle1"
              >
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
              {isLoggedIn && userRole !== "redactor" && (
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
  author: PropTypes.object.isRequired,
  onDelete: PropTypes.func.isRequired,
};

export default Comentario;
