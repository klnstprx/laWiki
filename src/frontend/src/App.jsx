import { Route, Routes } from "react-router-dom";
import HomePage from "./pages/HomePage.jsx";
import EntradaPage from "./pages/EntradaPage.jsx";
import WikiPage from "./pages/WikiPage.jsx";
import VersionPage from "./pages/VersionPage.jsx";
import MainLayout from "./layout/MainLayout.jsx";
import AdvancedSearchPage from "./pages/AdvancedSearchPage.jsx";
import FormWikiPage from "./pages/FormWikiPage.jsx";
import ProfilePage from "./pages/ProfilePage.jsx";
import FormEntradaPage from "./pages/FormEntradaPage.jsx";

function App() {
  return (
    <Routes>
      <Route element={<MainLayout />}>
        <Route path="/" element={<HomePage />} />
        <Route path="/entrada/:entryId/:versionId?" element={<EntradaPage />} />
        <Route path="/entrada/form/:wikiId" element={<FormEntradaPage />} />
        <Route
          path="/entrada/form/:entryId/:versionId?"
          element={<FormEntradaPage />}
        />
        <Route path="/wiki/:wikiId" element={<WikiPage />} />
        <Route path="/wiki/form/:wikiId" element={<FormWikiPage />} />
        <Route path="/wiki/form" element={<FormWikiPage />} />
        <Route path="/entrada/:entryId/:versionId?" element={<EntradaPage />} />
        <Route path="/versiones/:entryId" element={<VersionPage />} />
        <Route path="/advanced-search" element={<AdvancedSearchPage />} />
        <Route path="/perfil/:id" element={<ProfilePage />} />
      </Route>
    </Routes>
  );
}

export default App;
