import PropTypes from "prop-types";
import { Typography, Box } from "@mui/material";
import DOMPurify from "dompurify";

const Version = ({ content, editor, created_at }) => {
  return (
    <Box>
      <Typography variant="body2" color="text.secondary" gutterBottom>
        Editor: {editor}
      </Typography>
      <Typography variant="body2" color="text.secondary" gutterBottom>
        Fecha de creación: {new Date(created_at).toLocaleDateString()}
      </Typography>
      <Box
        sx={{ mt: 2 }}
        dangerouslySetInnerHTML={{
          __html: DOMPurify.sanitize(content),
        }}
      />
    </Box>
  );
};

Version.propTypes = {
  content: PropTypes.string.isRequired,
  editor: PropTypes.string.isRequired,
  created_at: PropTypes.string.isRequired,
};

export default Version;
