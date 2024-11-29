import { useEffect, useState } from "react";
import {
  getAllComentariosByVersion,
  getEntrada,
  getVersionById,
} from "../api.js";
import { useNavigate, useSearchParams } from "react-router-dom";
import MainLayout from "../layout/MainLayout.jsx";
import { useToast } from "../context/ToastContext.jsx";




function VersionPage() {

  const [version, setVersion] = useState({});
  const [error, setError] = useState(null);

  const [searchParams] = useSearchParams();
  const id = searchParams.get("id");



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
          <h1 style={{ fontSize: "36px", margin: "0" }}>Versiones Entrada</h1>
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

        <section
          style={{
            padding: "30px",
            backgroundColor: "white",
            margin: "20px 0",
            borderRadius: "8px",
          }}
        >
          
          {/* aqui van las versiones */}





        </section>
      </div>

    </MainLayout>
    
  );
}

export default VersionPage;
