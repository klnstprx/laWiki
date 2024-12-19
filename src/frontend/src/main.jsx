import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import App from "./App.jsx";
import "@fontsource/roboto/300.css";
import "@fontsource/roboto/400.css";
import "@fontsource/roboto/500.css";
import "@fontsource/roboto/700.css";
import { BrowserRouter } from "react-router-dom";
import { ToastProvider } from "./context/ToastProvider";
import { ThemeProvider } from "@mui/material/styles";
import { LocalizationProvider } from "@mui/x-date-pickers/LocalizationProvider";
import { AdapterDayjs } from "@mui/x-date-pickers/AdapterDayjs";
import theme from "./styles/theme";
import { GoogleOAuthProvider } from "@react-oauth/google";

const clientId = "292807851260-m45leggn30hmsvp8vhko14ue0lfi5276.apps.googleusercontent.com";

createRoot(document.getElementById("root")).render(
  <StrictMode>
    <GoogleOAuthProvider clientId={clientId}>
      <ThemeProvider theme={theme}>
        <BrowserRouter>
          <LocalizationProvider dateAdapter={AdapterDayjs}>
            <ToastProvider>
              <App />
            </ToastProvider>
          </LocalizationProvider>
        </BrowserRouter>
      </ThemeProvider>
    </GoogleOAuthProvider>
  </StrictMode>,
);
