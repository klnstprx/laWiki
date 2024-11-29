import { NavLink } from "react-router-dom";

const Footer = () => {
  return (
    <footer className="custom-footer">
      <p className="m-0">
        &copy; 2024 LaWiki. Todos los derechos reservados.{" "}
        <NavLink to="/" style={{ color: "inherit" }}>
          Home
        </NavLink>
      </p>
    </footer>
  );
};

export default Footer;