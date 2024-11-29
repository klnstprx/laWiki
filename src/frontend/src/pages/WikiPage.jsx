import { useEffect, useState } from "react";
import List from "@mui/material/List";
import ListItem from "@mui/material/ListItem";
import Alert from "@mui/material/Alert";
import Typography from "@mui/material/Typography";
import { getEntradasByWikiId, getWiki } from "../api.js";
import { useNavigate, useSearchParams } from "react-router-dom";
import EntradaCard from "../components/EntradaCard.jsx";
import { useToast } from "../context/ToastContext.jsx";
import MainLayout from "../layout/MainLayout.jsx";
import ConfirmationModal from "../components/ConfirmationModal.jsx";

function WikiPage() {
  const [wiki, setWiki] = useState({});
  const [entradas, setEntradas] = useState([]);
  const [error, setError] = useState(null);

  const [searchParams] = useSearchParams();
  const id = searchParams.get("id");
  const wikiId = searchParams.get("wikiID");

  const [showModal, setShowModal] = useState(false);
  const navigate = useNavigate();
  const { showToast } = useToast();

  const handleClose = () => {
    setShowModal(false);
    showToast("La entrada no se ha creado", "bg-danger");
  };

  const handleConfirm = () => {
    setShowModal(false);
    showToast("La entrada se ha creado correctamente!", "bg-success");
    navigate(`/wiki/${wikiId}`);
  };

  useEffect(() => {
    getWiki(id)
      .then(setWiki)
      .catch((err) => setError(err.message));
  }, [id]);

  useEffect(() => {
    getEntradasByWikiId(wikiId)
      .then(setEntradas)
      .catch((err) => setError(err.message));
  }, [wikiId]);

  async function enviarJSON(event) {
    console.log("Enviando formulario...");
    // Prevenir el envío normal del formulario
    event.preventDefault();

    // Obtener los datos del formulario
    const form = event.target; // El formulario
    const formData = new FormData(form); // Recoge todos los campos

    // Convertir FormData de wiki a un objeto JSON
    const jsonData = {};
    formData.forEach((value, key) => {
      jsonData[key] = value;
    });

    // Hacer la solicitud POST
    try {
      const response = await fetch("http://localhost:8000/api/entries/", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(jsonData),
      });

      // Manejar la respuesta
      if (response.ok) {
        const result = await response.json();
        console.log("Respuesta del servidor:", result);
      } else {
        console.error("Error en la respuesta:", response.status);
      }
    } catch (error) {
      console.error("Error al enviar:", error);
    }
  }

  {/*La URL es de este tipo http://localhost:5173/wiki?id=67311bf03399f3b49ccb8072&wikiId=67311c0143d96ecd81728a94 */ }

  return (
    <MainLayout>
      <div
        style={{
          fontFamily: "'Arial', sans-serif",
          backgroundColor: "#f5f5f5",
          padding: "40px",
          maxWidth: "1200px", // Aumentamos el ancho máximo para pantallas grandes
          margin: "0 auto",
          border: "1px solid #ddd",
          borderRadius: "8px",
          boxShadow: "0 4px 12px rgba(0, 0, 0, 0.1)", // Sombra más suave para el contenedor
          color: "black",
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
          {entradas.length > 0
            ? (
              <List>
                {entradas.map((comentario) => (
                  <ListItem
                    key={entradas.id}
                    style={{
                      borderBottom: "1px solid #ddd",
                      padding: "15px 0",
                    }}
                  >
                    <EntradaCard
                      title={comentario.title}
                      author={comentario.author}
                      created_at={comentario.created_at}
                    />
                  </ListItem>
                ))}
              </List>
            )
            : <Alert>No entries found.</Alert>}
        </section>

        {/* Formulario para añadir entradas */}
        <section
          style={{
            padding: "30px",
            backgroundColor: "white",
            border: "1px solid #ddd",
            borderRadius: "8px",
            boxShadow: "0 4px 8px rgba(0, 0, 0, 0.1)",
          }}
        >
          <h2 style={{ fontSize: "28px", marginBottom: "20px" }}>
            Añadir Entrada
          </h2>

          <form id="miFormulario" onSubmit={enviarJSON}>
            <div style={{ marginBottom: "20px" }}>
              <label
                htmlFor="author"
                style={{ fontWeight: "bold", fontSize: "18px" }}
              >
                Título:
              </label>
              <input
                type="text"
                id="title"
                name="title"
                required
                style={{
                  width: "100%",
                  padding: "12px",
                  marginTop: "10px",
                  border: "1px solid #ddd",
                  borderRadius: "6px",
                  fontSize: "16px",
                }}
              />
            </div>

            <div style={{ marginBottom: "20px" }}>
              <label
                htmlFor="author"
                style={{ fontWeight: "bold", fontSize: "18px" }}
              >
                Autor:
              </label>
              <input
                type="text"
                id="author"
                name="author"
                required
                style={{
                  width: "100%",
                  padding: "12px",
                  marginTop: "10px",
                  border: "1px solid #ddd",
                  borderRadius: "6px",
                  fontSize: "16px",
                }}
              />
            </div>

            <button
              type="submit"
              style={{
                width: "100%",
                padding: "15px",
                fontSize: "18px",
                backgroundColor: "#3c4f76",
                color: "white",
                border: "none",
                borderRadius: "6px",
                cursor: "pointer",
                boxShadow: "0 4px 8px rgba(0, 0, 0, 0.2)",
              }}
            >
              Enviar
            </button>
          </form>
        </section>
      </div>

      <ConfirmationModal
        message="¿Estás seguro de que quieres crear esta entrada?"
        show={showModal}
        handleClose={handleClose}
        handleConfirm={handleConfirm}
      >
      </ConfirmationModal>
    </MainLayout>
  );
}

export default WikiPage;
