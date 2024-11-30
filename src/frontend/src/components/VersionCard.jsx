import PropTypes from "prop-types";
import { Link } from "react-router-dom";
import { Card, CardContent, Typography, Grid2 } from "@mui/material";

const VersionCard = ({ entradaId, versionId, editor, created_at }) => {
  return (
    <Card sx={{ mb: 2 }}>
      <CardContent>
        <Grid2 container spacing={2} alignItems="center">
          <Grid2 xs={12} sm={5}>
            <Typography variant="body1">
              <strong>Fecha:</strong>{" "}
              {new Date(created_at).toLocaleDateString()}
            </Typography>
          </Grid2>
          <Grid2 xs={12} sm={5}>
            <Typography variant="body1">
              <strong>Editor:</strong> {editor}
            </Typography>
          </Grid2>
          <Grid2 xs={12} sm={2}>
            <Typography variant="body1">
              <Link to={`/entrada/${entradaId}/${versionId}`}>Ver</Link>
            </Typography>
          </Grid2>
        </Grid2>
      </CardContent>
    </Card>
  );
};

VersionCard.propTypes = {
  entradaId: PropTypes.string.isRequired,
  versionId: PropTypes.string.isRequired,
  editor: PropTypes.string.isRequired,
  created_at: PropTypes.string.isRequired,
};

export default VersionCard;
