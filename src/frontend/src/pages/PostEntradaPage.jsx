import { useState } from "react";
import { postEntry } from "../api/EntryApi.js";
import Button from "@mui/material/Button";
import { useNavigate, useSearchParams, Link } from "react-router-dom";
//import { useToast } from "../context/ToastContext.jsx";
import MainLayout from "../layout/MainLayout.jsx";

function PostEntradaPage() {
  const [error, setError] = useState(null);
  const navigate = useNavigate();
  //const { showToast } = useToast();
  const [searchParams] = useSearchParams();
  const id = searchParams.get("id");

  async function enviarJSON(event) {
    event.preventDefault();

    const form = event.target;
    const formData = new FormData(form);

    const jsonData = {};
    formData.forEach((value, key) => {
      jsonData[key] = value;
    });

    // Añadir el wiki_id automáticamente
    jsonData["wiki_id"] = id;

    try {
      const result = await postEntry(jsonData);
      console.log("Respuesta del servidor:", result);

      //showToast("Entrada creada correctamente", "bg-success");
      navigate(`/wiki?id=${id}`);
    } catch (error) {
      setError("Error al crear la entrada");
      console.error("Error al enviar:", error);
    }
  }

  return (
    <MainLayout>
      <div
        style={{
          padding: "30px",
          backgroundColor: "white",
          margin: "40px auto",
          maxWidth: "600px",
          border: "1px solid #ddd",
          borderRadius: "8px",
          boxShadow: "0 4px 12px rgba(0, 0, 0, 0.1)",
        }}
      >
        <h2 style={{ fontSize: "28px", marginBottom: "20px", textAlign: "center", color: "black" }}>
          Crear Nueva Entrada
        </h2>

        {error && (
          <div
            style={{
              backgroundColor: "#e57373",
              padding: "10px",
              borderRadius: "4px",
              marginBottom: "20px",
              textAlign: "center",
            }}
          >
            {error}
          </div>
        )}

        <form onSubmit={enviarJSON}>
          <div style={{ marginBottom: "20px" }}>
            <label
              htmlFor="title"
              style={{ fontWeight: "bold", fontSize: "18px", color: "black" }}
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
              style={{ fontWeight: "bold", fontSize: "18px", color: "black" }}
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
            Crear Entrada
          </button>

          <Link to={`/wiki?id=${id}`} style={{ textDecoration: "none" }}>
            <Button
              variant="contained"
              color="primary"
              style={{ marginTop: "20px" }}
            >
              Cancelar
            </Button>
          </Link>
        </form>
      </div>
    </MainLayout>
  );
}

export default PostEntradaPage;
