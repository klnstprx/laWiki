import { Box, CssBaseline } from "@mui/material";
import Header from "../components/Header";
import Footer from "../components/Footer";
import { Outlet } from "react-router-dom";

function MainLayout() {
  return (
    <Box sx={{ display: "flex", flexDirection: "column", minHeight: "100vh" }}>
      <CssBaseline />
      <Header />
      <Box sx={{ flexGrow: 1, pb: 5 }}>
        <Outlet />
      </Box>
      <Footer />
    </Box>
  );
}

export default MainLayout;
