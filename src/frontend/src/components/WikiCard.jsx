import PropTypes from "prop-types";
import {
  Card,
  CardActionArea,
  CardMedia,
  CardContent,
  Typography,
} from "@mui/material";
import { Link } from "react-router-dom";

const WikiCard = ({ wiki }) => {
  return (
    <Card sx={{ width: "100%" }}>
      <CardActionArea component={Link} to={`/wiki/${wiki.id}`}>
        <CardMedia
          component="img"
          height="140"
          image="https://res.cloudinary.com/dxj6khc6b/image/upload/v1733081078/books.png" //later we will use wiki.image, queried from the backend
          alt="Imagen de la Wiki"
        />
        <CardContent>
          <Typography gutterBottom variant="h5" component="div">
            {wiki.title}
          </Typography>
          <Typography variant="body2" color="textSecondary">
            {wiki.description}
          </Typography>
        </CardContent>
      </CardActionArea>
    </Card>
  );
};
WikiCard.propTypes = {
  wiki: PropTypes.shape({
    id: PropTypes.string.isRequired,
    title: PropTypes.string.isRequired,
    description: PropTypes.string.isRequired,
    image: PropTypes.string,
  }).isRequired,
};

export default WikiCard;
