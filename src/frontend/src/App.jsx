import { Link, Route, Routes } from "react-router-dom";
import HomePage from "./pages/HomePage.jsx";
import EntradaPage from "./pages/EntradaPage.jsx";

function App() {
  return (
    <div>
      <nav>
        <Link to="/">Home</Link>
      </nav>
      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/entrada" element={<EntradaPage />} />
        {/* <Route path="/about" element={<About />} /> */}
        {/* <Route path="/articles/:id" element={<Article />} /> */}
      </Routes>
    </div>
  );
}

export default App;
