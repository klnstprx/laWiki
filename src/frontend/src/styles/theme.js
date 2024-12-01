import { createTheme } from "@mui/material/styles";

const theme = createTheme({
  palette: {
    mode: "light",
    primary: {
      main: "#52796f",
      light: "#79a28f",
      dark: "#2f5248",
      contrastText: "#ffffff",
    },
    secondary: {
      main: "#C84630",
      light: "#d9705a",
      dark: "#912b1f",
      contrastText: "#ffffff",
    },
    error: {
      main: "#bc252a",
    },
    warning: {
      main: "#FB5607",
    },
    success: {
      main: "#6ea683",
    },
    background: {
      default: "#F8F9FA",
      paper: "#FFFFFF",
    },
  },
  components: {
    // Style overrides for buttons
    MuiButton: {
      styleOverrides: {
        root: {
          textTransform: "none", // Disable uppercase transformation
          borderRadius: 8, // Rounded corners
        },
        containedPrimary: {
          "&:hover": {
            backgroundColor: "#2f5248", // Darken primary color on hover
          },
        },
        containedSecondary: {
          "&:hover": {
            backgroundColor: "#912b1f", // Darken secondary color on hover
          },
        },
      },
    },
    // Style overrides for other components if needed
    MuiListItem: {
      styleOverrides: {
        root: {
          "&.Mui-selected": {
            backgroundColor: "#79a28f", // Highlight selected items
            "&:hover": {
              backgroundColor: "#6a9481", // Darken on hover when selected
            },
          },
          "&:hover": {
            backgroundColor: "#f0f0f0", // Light grey on hover
          },
        },
      },
    },
  },
});

export default theme;
