import { Link, Route, Routes } from "react-router-dom";
import Home from "./pages/Home.jsx";

function App() {
  return (
    <div>
      <nav>
        <Link to="/">Home</Link>
      </nav>
      <Routes>
        <Route path="/" element={<Home />} />
        {/* <Route path="/about" element={<About />} /> */}
        {/* <Route path="/articles/:id" element={<Article />} /> */}
      </Routes>
    </div>
  );
}

export default App;
