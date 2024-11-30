import { Avatar, Typography, IconButton, Box } from "@mui/material";
import DeleteIcon from "@mui/icons-material/Delete";
import PropTypes from "prop-types";

const Comentario = ({ id, content, rating, created_at, author, onDelete }) => {
  const handleDelete = () => {
    onDelete(id);
  };

  return (
    <Box display="flex" mb={2}>
      <Avatar
        src={`https://ui-avatars.com/api/?name=${author}&background=random`}
        alt={author}
        sx={{ mr: 2 }}
      />
      <Box flexGrow={1}>
        <Box display="flex" justifyContent="space-between">
          <Typography variant="subtitle1">{author}</Typography>
          <Typography variant="caption" color="textSecondary">
            {new Date(created_at).toLocaleDateString()}
          </Typography>
        </Box>
        <Typography variant="body1" color="textPrimary" gutterBottom>
          {content}
        </Typography>
        <Box display="flex" justifyContent="space-between" alignItems="center">
          <Typography variant="body2" color="textSecondary">
            Rating: {rating}/5
          </Typography>
          <IconButton color="error" onClick={handleDelete}>
            <DeleteIcon />
          </IconButton>
        </Box>
      </Box>
    </Box>
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
