import { Route, Routes } from "react-router-dom";
import HomePage from "./pages/HomePage.jsx";
import EntradaPage from "./pages/EntradaPage.jsx";
import WikiPage from "./pages/WikiPage.jsx";
import VersionPage from "./pages/VersionPage.jsx";
import FormEntradaPage from "./pages/FormEntradaPage.jsx";
import MainLayout from "./layout/MainLayout.jsx";
import EditarEntradaPage from "./pages/EditarEntradaPage.jsx";
import AdvancedSearchPage from "./pages/AdvancedSearchPage.jsx";
import FormWikiPage from "./pages/FormWikiPage.jsx";

function App() {
  return (
    <Routes>
      <Route element={<MainLayout />}>
        <Route path="/" element={<HomePage />} />
        <Route path="/entrada/:entryId/:versionId?" element={<EntradaPage />} />
        <Route path="/wiki/:id" element={<WikiPage />} />
        <Route path="/wiki/form" element={<FormWikiPage />} />
        <Route path="/wiki/form/:wikiId" element={<FormWikiPage />} />
        <Route path="/versiones/:entryId" element={<VersionPage />} />
        <Route path="/crear-entrada/:id" element={<FormEntradaPage />} />
        <Route
          path="/editarEntrada/:entryId/:versionId?"
          element={<EditarEntradaPage />}
        />
        <Route path="/advanced-search" element={<AdvancedSearchPage />} />
        {/* <Route path="/about" element={<About />} /> */}
        {/* <Route path="/articles/:id" element={<Article />} /> */}
      </Route>
    </Routes>
  );
}

export default App;
