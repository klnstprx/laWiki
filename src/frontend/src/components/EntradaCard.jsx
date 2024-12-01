import { Link } from "react-router-dom";
import PropTypes from "prop-types";
import {
  Card,
  CardActionArea,
  CardContent,
  Typography,
  Grid2,
  IconButton
} from "@mui/material";
import DeleteIcon from "@mui/icons-material/Delete";

const EntradaCard = ({ id, title, author, createdAt, onEntradaClick, onDelete }) => {
  const handleClick = () => {
    if (onEntradaClick) {
      onEntradaClick(id);
    }
  };

  const handleDelete = () => {
    onDelete(id);
  }; 

  return (
    <Card sx={{ mb: 2 }}>
      <CardActionArea
        component={Link}
        to={`/entrada/${id}`}
        onClick={handleClick}
      >
        <CardContent>
          <Typography variant="h5" component="div">
            {title}
          </Typography>
          <Grid2 container spacing={2} sx={{ mt: 1 }}>
            <Grid2 xs={6}>
              <Typography variant="subtitle1" color="textSecondary">
                Author
              </Typography>
              <Typography variant="body2">{author}</Typography>
            </Grid2>
            <Grid2 xs={6}>
              <Typography variant="subtitle1" color="textSecondary">
                Created At
              </Typography>
              <Typography variant="body2">{createdAt}</Typography>
            </Grid2>
            <Grid2>
              <IconButton color="error" onClick={handleDelete}>
                  <DeleteIcon />
              </IconButton>
            </Grid2>
          </Grid2>
        </CardContent>
      </CardActionArea>
    </Card>
  );
};

EntradaCard.propTypes = {
  id: PropTypes.string.isRequired,
  title: PropTypes.string.isRequired,
  author: PropTypes.string.isRequired,
  createdAt: PropTypes.string.isRequired,
  onEntradaClick: PropTypes.func,
  onDelete: PropTypes.func.isRequired,
};

export default EntradaCard;
