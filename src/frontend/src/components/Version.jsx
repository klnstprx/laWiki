import Paper from "@mui/material/Paper";
import { MapContainer, TileLayer, Marker, Popup } from "react-leaflet";
import PropTypes from "prop-types";
import { Typography, Box } from "@mui/material";
import DOMPurify from "dompurify";
import "leaflet/dist/leaflet.css";

const Version = ({ content, editor, created_at, address, coordinates }) => {
  return (
    <div>
      <Typography variant="h6" gutterBottom>
        Editor: {editor}
      </Typography>
      <Typography variant="h6" gutterBottom>
        Fecha de creación: {created_at}
      </Typography>
      <Typography variant="h6" gutterBottom>
        Contenido:
      </Typography>
      <Box
        sx={{ mt: 2 }}
        dangerouslySetInnerHTML={{
          __html: DOMPurify.sanitize(content),
        }}
      ></Box>
      <Box variant="h6" gutterBottom>
        Ubicación: {address || "No especificada"}
      </Box>

      {coordinates && (
        <Paper elevation={3} sx={{ mt: 3 }}>
          <Typography variant="body2">
            Coordenadas: Lat: {coordinates.lat}, Lon: {coordinates.lon}
          </Typography>
          <Typography variant="h6" gutterBottom>
            Ubicación
          </Typography>
          <MapContainer
            center={[coordinates.lat, coordinates.lon]}
            zoom={15}
            style={{ height: "400px", width: "100%" }}
          >
            <TileLayer
              url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
              attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
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
};

export default Version;
