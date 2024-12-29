import { useEffect, useState } from "react";
import Paper from "@mui/material/Paper";
import { MapContainer, Marker, Popup, TileLayer } from "react-leaflet";
import markerIconPng from "leaflet/dist/images/marker-icon.png";
import { Icon } from "leaflet";
import PropTypes from "prop-types";
import {
  Avatar,
  Box,
  CircularProgress,
  Divider,
  Stack,
  Typography,
} from "@mui/material";
import { useTheme } from "@mui/material/styles";
import DOMPurify from "dompurify";
import "leaflet/dist/leaflet.css";
import Slider from "react-slick";
import "slick-carousel/slick/slick.css";
import "slick-carousel/slick/slick-theme.css";
import { getMedia } from "../api/MediaApi";

import ArrowForwardIosIcon from "@mui/icons-material/ArrowForwardIos";
import Dialog from "@mui/material/Dialog";
import DialogContent from "@mui/material/DialogContent";
import IconButton from "@mui/material/IconButton";
import CloseIcon from "@mui/icons-material/Close";
import { getUser } from "../api/AuthApi";

// Custom Previous Arrow Component
function PrevArrow(props) {
  const { onClick } = props;
  const theme = useTheme();

  return (
    <div onClick={onClick}>
      <ArrowForwardIosIcon
        fontSize="large"
        style={{ transform: "rotate(180deg)" }}
        sx={{
          left: 0,
          zIndex: 1,
          position: "absolute",
          top: "50%",
          transform: "translateY(-50%)",
          color: theme.palette.primary.dark,
          "&:hover": {
            color: theme.palette.primary.main,
          },
        }}
      />
    </div>
  );
}

// Custom Next Arrow Component
function NextArrow(props) {
  const { onClick } = props;
  const theme = useTheme();

  return (
    <div onClick={onClick}>
      <ArrowForwardIosIcon
        fontSize="large"
        sx={{
          right: 0,
          zIndex: 1,
          position: "absolute",
          top: "50%",
          transform: "translateY(-50%)",
          color: theme.palette.primary.dark,
          "&:hover": {
            color: theme.palette.primary.main,
          },
        }}
      />
    </div>
  );
}

const Version = ({
  content,
  editor,
  created_at,
  address,
  coordinates,
  media_ids,
}) => {
  const [medias, setMedias] = useState(null);
  const [mediaError, setMediaError] = useState(null);
  const [loading, setLoading] = useState(true);
  const [usuario, setUsuario] = useState({}); // add state

  // New state variables for dialog
  const [openDialog, setOpenDialog] = useState(false);
  const [selectedImage, setSelectedImage] = useState(null);

  //carga el usuario

//carga el usuario
  useEffect(() => {
    const fetchUsuario = async () => {
      try {
        const userData = await getUser(editor);
        setUsuario(userData);
      } catch (error) {
        console.error("Error fetching user:", error);
      }
    };

    if (editor) {
      fetchUsuario();
    }
  }, [editor]);

  useEffect(() => {
    const fetchMedia = async () => {
      try {
        const mediasData = await Promise.all(
          media_ids.map((id) => getMedia(id)),
        );
        setMedias(mediasData);
      } catch (error) {
        console.error("Error fetching media:", error);
        setMediaError("Failed to load images.");
      } finally {
        setLoading(false);
      }
    };

    if (media_ids && media_ids.length > 0) {
      fetchMedia();
    } else {
      setLoading(false);
    }
  }, [media_ids]);

  // Slider settings for react-slick
  const sliderSettings = {
    arrows: true,
    dots: true,
    adaptiveHeight: true,
    slidesToShow: 1,
    slidesToScroll: 1,
    infinite: true,
    nextArrow: <NextArrow />,
    prevArrow: <PrevArrow />,
  };

  return (
    <div>
      {/* Header Section */}
      <Stack direction="row" justifyContent="space-between" alignItems="center">
        <Box sx={{ display: "flex", alignItems: "center", gap: "1em" }}>
          <Avatar
            src={usuario.picture}
            alt={editor}
          />
          <Typography variant="subtitle1" fontWeight="bold">
            <a href={`/perfil/${usuario.id}`}>{usuario.name}</a>
          </Typography>
        </Box>
        <Typography variant="caption" color="text.secondary">
          {new Date(created_at).toLocaleDateString()}
        </Typography>
      </Stack>

      <Divider sx={{ my: 2 }} />

      {/* Content Section */}
      <Box
        sx={{ mt: 2 }}
        dangerouslySetInnerHTML={{
          __html: DOMPurify.sanitize(content),
        }}
      >
      </Box>

      <Divider sx={{ my: 2 }} />

      {/* Media Section */}
      {loading
        ? (
          <Box sx={{ display: "flex", justifyContent: "center", mt: 3 }}>
            <CircularProgress />
          </Box>
        )
        : mediaError
        ? (
          <Typography color="error" sx={{ mt: 3 }}>
            {mediaError}
          </Typography>
        )
        : medias && medias.length > 0
        ? (
          <>
            <Slider {...sliderSettings}>
              {medias.map((media, index) => (
                <div key={index}>
                  <Box sx={{ m: 3 }}>
                    <img
                      src={media.uploadUrl}
                      alt={`Image ${index + 1}`}
                      style={{
                        maxWidth: "100%",
                        maxHeight: "500px",
                        width: "auto",
                        height: "auto",
                        display: "block",
                        margin: "10px auto",
                        borderRadius: "8px",
                        boxShadow: "0px 0px 10px 0px rgba(0,0,0,0.75)",
                        cursor: "pointer",
                      }}
                      onClick={() => {
                        setSelectedImage(media.uploadUrl);
                        setOpenDialog(true);
                      }}
                    />
                  </Box>
                </div>
              ))}
            </Slider>

            {/* Image Dialog */}
            <Dialog
              open={openDialog}
              onClose={() => setOpenDialog(false)}
              maxWidth="xl"
              sx={{
                "& .MuiDialog-paper": {
                  maxWidth: "80%",
                },
              }}
            >
              <DialogContent sx={{ padding: 0, position: "relative" }}>
                <IconButton
                  aria-label="close"
                  onClick={() => setOpenDialog(false)}
                  sx={{
                    position: "absolute",
                    right: 8,
                    top: 8,
                    color: (theme) => theme.palette.grey[500],
                    zIndex: 1,
                  }}
                >
                  <CloseIcon />
                </IconButton>
                <img
                  src={selectedImage}
                  alt="Selected"
                  style={{
                    width: "100%",
                    height: "100%",
                    objectFit: "contain",
                    display: "block",
                  }}
                />
              </DialogContent>
            </Dialog>
          </>
        )
        : null}

      {/* Map Section */}
      {coordinates && (
        <Paper elevation={3} sx={{ mt: 3 }}>
          <Typography variant="body2" sx={{ p: 2 }}>
            Coordinates: Lat: {coordinates.lat}, Lon: {coordinates.lon}
          </Typography>
          <Typography variant="body2" sx={{ p: 2 }}>
            Location: {address || "Not specified"}
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
            <Marker
              position={[coordinates.lat, coordinates.lon]}
              icon={new Icon({
                iconUrl: markerIconPng,
                iconSize: [25, 41],
                iconAnchor: [12, 41],
              })}
            >
              <Popup>Location associated with this version</Popup>
            </Marker>
          </MapContainer>
        </Paper>
      )}

      {!coordinates && address && (
        <Typography variant="body1" sx={{ mt: 2 }}>
          Unable to retrieve location for the provided address.
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
