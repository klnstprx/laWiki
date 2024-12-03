import { useState, useEffect } from "react";
import Paper from "@mui/material/Paper";
import { MapContainer, TileLayer, Marker, Popup } from "react-leaflet";
import PropTypes from "prop-types";
import { Typography, Box, Divider, Stack, CardMedia } from "@mui/material";
import DOMPurify from "dompurify";
import "leaflet/dist/leaflet.css";
import "react-responsive-carousel/lib/styles/carousel.min.css";
import { Carousel } from "react-responsive-carousel";
import { getMedia } from "../api/MediaApi";

const Version = ({ content, editor, created_at, address, coordinates, media_ids }) => {
  const [medias, setMedias] = useState(null);
  const [mediaError, setMediaError] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchMedia = async () => {
      try {
        const mediasData = await Promise.all(media_ids.map((id) => getMedia(id)));
        setMedias(mediasData);
      } catch (error) {
        console.error("Error fetching media:", error);
        setMediaError("Failed to load image.");
      } finally {
        setLoading(false);
      }
    };

    if (media_ids) {
      fetchMedia();
    } else {
      setLoading(false);
    }
  }, [media_ids]);

  return (
    <div>
      {/* Content Section */}
      <Box
        sx={{ mt: 2 }}
        dangerouslySetInnerHTML={{
          __html: DOMPurify.sanitize(content),
        }}
      ></Box>

      <Divider sx={{ my: 4 }} />

      {/* Details Section */}
      <Stack
        direction="row"
        spacing={2}
        alignItems="center"
        justifyContent="space-between"
        sx={{ mx: "auto", width: "80%" }}
      >
        <Typography variant="subtitle1" gutterBottom>
          Editor: {editor}
        </Typography>
        <Divider orientation="vertical" flexItem />
        <Typography variant="subtitle1" gutterBottom>
          Fecha de creación:{" "}
          {new Date(created_at).toLocaleString("es-ES", {
            year: "numeric",
            month: "short",
            day: "numeric",
            hour: "2-digit",
            minute: "2-digit",
          })}
        </Typography>
        <Divider orientation="vertical" flexItem />
        <Typography variant="subtitle1" gutterBottom>
          Ubicación: {address || "No especificada"}
        </Typography>
      </Stack>

      {// show a carousel with images if there are any}
      medias && medias.length > 0 && (
        <Paper elevation={3} sx={{ mt: 3 }}>
          <Carousel showArrows={true} showThumbs={false}>
            {medias.map((media) => (
              <div key={media.id}>
                <CardMedia
                  component="img"
                  image={media.uploadUrl}
                  alt="Imagen de la Wiki"
                />
              </div>
            ))}
          </Carousel>
        </Paper>
      )}

      {/* Map Section */}
      {coordinates && (
        <Paper elevation={3} sx={{ mt: 3 }}>
          <Typography variant="body2" sx={{ p: 2 }}>
            Coordenadas: Lat: {coordinates.lat}, Lon: {coordinates.lon}
          </Typography>
          <MapContainer
            center={[coordinates.lat, coordinates.lon]}
            zoom={15}
            style={{ height: "400px", width: "100%" }}
          >
            <TileLayer
              url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
              attribution="&copy; OpenStreetMap contributors"
            />
            <Marker position={[coordinates.lat, coordinates.lon]}>
              <Popup>Ubicación asociada a esta versión</Popup>
            </Marker>
          </MapContainer>
        </Paper>
      )}

      {!coordinates && address && (
        <Typography variant="body1" sx={{ mt: 2 }}>
          No se pudo obtener la ubicación para la dirección proporcionada.
        </Typography>
      )}
    </div>
  );
};

Version.propTypes = {
  content: PropTypes.string.isRequired,
  editor: PropTypes.string.isRequired,
  created_at: PropTypes.string.isRequired,
  address: PropTypes.string,
  coordinates: PropTypes.shape({
    lat: PropTypes.number,
    lon: PropTypes.number,
  }),
  media_ids: PropTypes.arrayOf(PropTypes.string),
};

export default Version;
