import { Route, Routes } from "react-router-dom";
import HomePage from "./pages/HomePage.jsx";
import EntradaPage from "./pages/EntradaPage.jsx";
import WikiPage from "./pages/WikiPage.jsx";
import VersionPage from "./pages/VersionPage.jsx";
import PostEntradaPage from "./pages/PostEntradaPage.jsx";
import MainLayout from "./layout/MainLayout.jsx";
import EditarEntradaPage from "./pages/EditarEntradaPage.jsx";

function App() {
  return (
    <Routes>
      <Route element={<MainLayout />}>
        <Route path="/" element={<HomePage />} />
        <Route path="/entrada/:entryId/:versionId?" element={<EntradaPage />} />
        <Route path="/wiki/:id" element={<WikiPage />} />
        <Route path="/versiones/:entryId" element={<VersionPage />} />
        <Route path="/crear-entrada" element={<PostEntradaPage />} />
        <Route path="/editarEntrada" element={<EditarEntradaPage />} />
        {/* <Route path="/about" element={<About />} /> */}
        {/* <Route path="/articles/:id" element={<Article />} /> */}
      </Route>
    </Routes>
  );
}

export default App;
