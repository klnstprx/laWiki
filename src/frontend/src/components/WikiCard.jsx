import { useState, useEffect } from "react";
import PropTypes from "prop-types";
import {
  Card,
  CardActionArea,
  CardMedia,
  CardContent,
  Typography,
} from "@mui/material";
import { Link } from "react-router-dom";
import { getMedia } from "../api/MediaApi";

const WikiCard = ({ wiki }) => {
  const [media, setMedia] = useState(null);
  const [mediaError, setMediaError] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchMedia = async () => {
      try {
        const mediaData = await getMedia(wiki.media_id);
        setMedia(mediaData);
      } catch (error) {
        console.error("Error fetching media:", error);
        setMediaError("Failed to load image.");
      } finally {
        setLoading(false);
      }
    };

    if (wiki.media_id) {
      fetchMedia();
    } else {
      setLoading(false);
    }
  }, [wiki.media_id]);

  if (loading) {
    return <Typography>Loading...</Typography>;
  }

  if (mediaError) {
    return <Typography color="error">{mediaError}</Typography>;
  }
  const defaultImageUrl = "https://res.cloudinary.com/dxj6khc6b/image/upload/v1733154871/abstract_background.png";
  return (
    <Card sx={{ width: "100%" }}>
      <CardActionArea component={Link} to={`/wiki/${wiki.id}`}>
        <CardMedia
          component="img"
          height="140"
          image={(media && media.uploadUrl) || defaultImageUrl}
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
    media_id: PropTypes.string,
  }).isRequired,
};

export default WikiCard;