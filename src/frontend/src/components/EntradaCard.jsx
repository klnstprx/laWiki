import { Link } from "react-router-dom";
import PropTypes from "prop-types";
import {
  Card,
  CardActionArea,
  CardContent,
  Typography,
  Grid2,
} from "@mui/material";

const EntradaCard = ({ id, title, author, createdAt, onEntradaClick }) => {
  const handleClick = () => {
    if (onEntradaClick) {
      onEntradaClick(id);
    }
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
};

export default EntradaCard;
