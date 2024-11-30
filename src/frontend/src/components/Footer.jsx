import { Box, Typography, Link as MuiLink } from "@mui/material";
import { Link } from "react-router-dom";

const Footer = () => {
  const currentYear = new Date().getFullYear();
  return (
    <Box
      component="footer"
      py={2}
      bgcolor="primary.main"
      color="white"
      textAlign="center"
      sx={{ position: "fixed", bottom: 0, left: 0, right: 0, zIndex: 1000 }}
      elevation={3}
    >
      <Typography variant="body2">
        &copy; {currentYear} LaWiki. Todos los derechos reservados.{" "}
        <MuiLink component={Link} to="/" color="inherit" underline="none">
          Home
        </MuiLink>
      </Typography>
    </Box>
  );
};

export default Footer;
