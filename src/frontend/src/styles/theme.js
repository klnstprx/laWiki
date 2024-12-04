import { createTheme } from "@mui/material/styles";

const theme = createTheme({
  palette: {
    mode: "light",
    primary: {
      main: "#52796f", // Teal Green
      light: "#84a29d", // Lighter Teal Green
      dark: "#36554a", // Darker Teal Green
      contrastText: "#ffffff",
    },
    secondary: {
      main: "#355c7d", // Muted Blue (Analogous to Teal Green)
      light: "#5d85a5", // Lighter Muted Blue
      dark: "#1e3b54", // Darker Muted Blue
      contrastText: "#ffffff",
    },
    error: {
      main: "#aa4465", // Muted Red-Purple (Complementary to Teal Green)
    },
    warning: {
      main: "#e7a977", // Muted Orange (Harmonious with Teal and Blue)
    },
    success: {
      main: "#6ea683", // Muted Green (Analogous to Teal Green)
    },
    background: {
      default: "#f5f7f6", // Very Light Teal Green Tint
      paper: "#ffffff",
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
            backgroundColor: "#36554a", // Darken primary color on hover
          },
        },
        containedSecondary: {
          "&:hover": {
            backgroundColor: "#1e3b54", // Darken secondary color on hover
          },
        },
      },
    },
    MuiListItem: {
      styleOverrides: {
        root: {
          "&.Mui-selected": {
            backgroundColor: "#84a29d", // Highlight selected items
            "&:hover": {
              backgroundColor: "#6d8e89", // Darken on hover when selected
            },
          },
          "&:hover": {
            backgroundColor: "#eef1f0", // Light grey on hover
          },
        },
      },
    },
    MuiBreadcrumbs: {
      styleOverrides: {
        root: {
          marginBottom: "16px",
        },
      },
    },
    MuiTypography: {
      styleOverrides: {
        root: {
          "&.breadcrumb-link": {
            textDecoration: "none",
            color: "#355c7d",
            fontWeight: 500,
            "&:hover": {
              color: "#52796f",
            },
          },
          "&.breadcrumb-active": {
            color: "#52796f",
            fontWeight: 500,
          },
        },
      },
    },
  },
});

export default theme;
