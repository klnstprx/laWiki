import React from 'react';
import { Card, CardActionArea, CardMedia, CardContent, Typography } from '@mui/material';
import { Link } from 'react-router-dom';

const WikiCard = ({ wiki }) => {
  return (
    <Card>
      <CardActionArea component={Link} to={`/wiki/${wiki.id}`}>
        <CardMedia
          component="img"
          height="140"
          image="https://via.placeholder.com/350x140" //later we will use wiki.image, queried from the backend
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

export default WikiCard;