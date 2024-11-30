import { useEffect, useState } from "react";
import List from "@mui/material/List";
import ListItem from "@mui/material/ListItem";
import Alert from "@mui/material/Alert";
import Button from "@mui/material/Button";
import Typography from "@mui/material/Typography";
import { searchEntries } from "../api/EntryApi.js";
import { getWiki } from "../api/WikiApi.js";
import { useSearchParams, Link } from "react-router-dom";
import EntradaCard from "../components/EntradaCard.jsx";
import MainLayout from "../layout/MainLayout.jsx";

function WikiPage() {
  const [wiki, setWiki] = useState({});
  const [entradas, setEntradas] = useState([]);
  const [error, setError] = useState(null);

  const [searchParams] = useSearchParams();
  const id = searchParams.get("id");

  useEffect(() => {
    getWiki(id)
      .then(setWiki)
      .catch((err) => setError(err.message));
  }, [id]);

  useEffect(() => {
    searchEntries({ wikiID: id })
      .then(setEntradas)
      .catch((err) => setError(err.message));
  }, [id]);

  return (
    <MainLayout>
      <div
        style={{
          fontFamily: "'Arial', sans-serif",
          backgroundColor: "#f5f5f5",
          padding: "40px",
          margin: "0 auto",
          border: "1px solid #ddd",
          borderRadius: "8px",
          boxShadow: "0 4px 12px rgba(0, 0, 0, 0.1)", // Sombra más suave para el contenedor
          color: "black",
          width: "100vw",
        }}
      >
        {/* Cabecera de la página */}
        <header
          style={{
            backgroundColor: "#3c4f76",
            color: "white",
            padding: "20px",
            borderRadius: "8px 8px 0 0",
            textAlign: "center",
          }}
        >
          <h1 style={{ fontSize: "36px", margin: "0" }}>Datos de la Wiki</h1>
          {error && (
            <div
              style={{
                backgroundColor: "#e57373",
                padding: "15px",
                marginTop: "15px",
                borderRadius: "4px",
              }}
            >
              <p>{error}</p>
            </div>
          )}
        </header>

        {/* Información de la wiki */}
        <section
          style={{
            padding: "30px",
            backgroundColor: "white",
            margin: "20px 0",
            borderRadius: "8px",
          }}
        >
          <Typography
            variant="h6"
            style={{ fontWeight: "bold", marginBottom: "10px" }}
          >
            Título: <span style={{ fontWeight: "normal" }}>{wiki.title}</span>
          </Typography>

          <Typography
            variant="h6"
            style={{ fontWeight: "bold", marginBottom: "10px" }}
          >
            Descripción:{" "}
            <span style={{ fontWeight: "normal" }}>{wiki.description}</span>
          </Typography>

          <Typography
            variant="h6"
            style={{ fontWeight: "bold", marginBottom: "10px" }}
          >
            Categoría:{" "}
            <span style={{ fontWeight: "normal" }}>{wiki.category}</span>
          </Typography>
        </section>

        {/* Entradas */}
        <section
          style={{
            padding: "30px",
            backgroundColor: "white",
            marginBottom: "20px",
            borderRadius: "8px",
          }}
        >
          <h2
            style={{
              fontSize: "28px",
              borderBottom: "2px solid #ddd",
              paddingBottom: "10px",
              marginBottom: "20px",
            }}
          >
            Entradas
          </h2>
          {entradas.length > 0 ? (
            <List>
              {entradas.map((entrada) => (
                <ListItem
                  key={entrada.id}
                  style={{
                    borderBottom: "1px solid #ddd",
                    padding: "15px 0",
                  }}
                >
                  <EntradaCard
                    id={entrada.id}
                    title={entrada.title}
                    author={entrada.author}
                    createdAt={entrada.created_at}
                  />
                </ListItem>
              ))}
            </List>
          ) : (
            <Alert>No entries found.</Alert>
          )}
        </section>
            
        {/* Botón para crear nueva entrada */}
        <Link to={`/crear-entrada?id=${id}`} style={{ textDecoration: "none" }}>
            <Button
                variant="contained"
                color="primary"
                style={{ marginTop: "20px" }}
            >
                Crear Nueva Entrada
            </Button>
        </Link>
      </div>
    </MainLayout>
  );
}

export default WikiPage;
