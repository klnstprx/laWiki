import { Route, Routes } from "react-router-dom";
import HomePage from "./pages/HomePage.jsx";
import EntradaPage from "./pages/EntradaPage.jsx";
import WikiPage from "./pages/WikiPage.jsx";
import VersionPage from "./pages/VersionPage.jsx";
import PostEntradaPage from "./pages/PostEntradaPage.jsx";

function App() {
  return (
    <div>
      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/entrada" element={<EntradaPage />} />
        <Route path="/wiki" element={<WikiPage />} />
        <Route path="/versiones" element={<VersionPage />} />
        <Route path="/crear-entrada" element={<PostEntradaPage />} />
        {/* <Route path="/about" element={<About />} /> */}
        {/* <Route path="/articles/:id" element={<Article />} /> */}
      </Routes>
    </div>
  );
}

export default App;
