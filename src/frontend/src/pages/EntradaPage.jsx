
import { useEffect, useState } from "react";
import List from "@mui/material/List";
import ListItem from "@mui/material/ListItem";
import Alert from "@mui/material/Alert";
import Typography from "@mui/material/Typography";
import { getEntrada, getAllComentariosByVersion, getVersionById } from "../api.js";
import { useSearchParams, useNavigate} from 'react-router-dom';
import ComentarioComponent from "../components/ComentarioComponent.jsx";
import VersionComponent from "../components/VersionComponent.jsx";


function EntradaPage() {
  
    const [entrada, setEntrada] = useState({});
    const [version, setVersion] = useState({});
    const [comentarios, setComentarios] = useState([]);
    const [error, setError] = useState(null);
    
    const [searchParams] = useSearchParams();

    const id = searchParams.get('id');
    const versionID = searchParams.get('versionID');

    useEffect(() => {
      getEntrada(id)
      .then(setEntrada)
      .catch((err) => setError(err.message));
  }, [id]);

  useEffect(() => {
      getAllComentariosByVersion(versionID)
        .then(setComentarios)
        .catch((err) => setError(err.message));
    }, [versionID]);

    useEffect(() => {
      getVersionById(versionID)
        .then(setVersion)
        .catch((err) => setError(err.message));
    }, [versionID]);



    async function enviarJSON(event) {
      console.log('Enviando formulario...');
      // Prevenir el envío normal del formulario
      event.preventDefault();

      // Obtener los datos del formulario
      const form = event.target; // El formulario
      const formData = new FormData(form); // Recoge todos los campos

      // Convertir FormData a un objeto JSON
      const jsonData = {};
      formData.forEach((value, key) => {
          jsonData[key] = key === "rating" ? parseInt(value, 10) : value;
      });

      jsonData['version_id'] = version.id
      // Hacer la solicitud POST
      try {
          const response = await fetch('http://localhost:8000/api/comments/', {
              method: 'POST',
              headers: {
                  'Content-Type': 'application/json',
              },
              body: JSON.stringify(jsonData),
          });

          // Manejar la respuesta
          if (response.ok) {
              const result = await response.json();
              console.log('Respuesta del servidor:', result);
          } else {
              console.error('Error en la respuesta:', response.status);
          }
      } catch (error) {
          console.error('Error al enviar:', error);
      }
    }

   

    {/*La URL es de este tipo http://localhost:5173/entrada?id=67311bf03399f3b49ccb8072&versionID=67311c0143d96ecd81728a94 */}

    return (
<div
  style={{
    fontFamily: "'Arial', sans-serif",
    backgroundColor: "#f5f5f5",
    padding: "40px",
    maxWidth: "1200px",  // Aumentamos el ancho máximo para pantallas grandes
    margin: "0 auto",
    border: "1px solid #ddd",
    borderRadius: "8px",
    boxShadow: "0 4px 12px rgba(0, 0, 0, 0.1)", // Sombra más suave para el contenedor
    color: "black"
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
      <div style={{ backgroundColor: "#e57373", padding: "15px", marginTop: "15px", borderRadius: "4px" }}>
        <p>{error}</p>
      </div>
    )}
  </header>

  {/* Información de la entrada */}
  <section style={{ padding: "30px", backgroundColor: "white", margin: "20px 0", borderRadius: "8px" }}>
    <Typography variant="h6" style={{ fontWeight: "bold", marginBottom: "10px" }}>
      Título: <span style={{ fontWeight: "normal"}}>{entrada.title}</span>
    </Typography>

    <Typography variant="h6" style={{ fontWeight: "bold", marginBottom: "10px" }}>
      Autor: <span style={{ fontWeight: "normal" }}>{entrada.author}</span>
    </Typography>

    <Typography variant="h6" style={{ fontWeight: "bold", marginBottom: "10px" }}>
      Fecha de creación: <span style={{ fontWeight: "normal" }}>{entrada.created_at}</span>
    </Typography>
  </section>

  {/* Contenido de la versión */}
  <section style={{ padding: "30px", backgroundColor: "white", marginBottom: "20px", borderRadius: "8px" }}>
    <h2 style={{ fontSize: "28px", borderBottom: "2px solid #ddd", paddingBottom: "10px", marginBottom: "20px" }}>
      Contenido de la versión
    </h2>
    <VersionComponent
      content={version.content}
      editor={version.editor}
      created_at={version.created_at}
      entry_id={version.entry_id}
    />
  </section>

  {/* Comentarios */}
  <section style={{ padding: "30px", backgroundColor: "white", marginBottom: "20px", borderRadius: "8px" }}>
    <h2 style={{ fontSize: "28px", borderBottom: "2px solid #ddd", paddingBottom: "10px", marginBottom: "20px" }}>
      Comentarios
    </h2>
    {comentarios.length > 0 ? (
      <List>
        {comentarios.map((comentario) => (
          <ListItem key={comentarios.id} style={{ borderBottom: "1px solid #ddd", padding: "15px 0" }}>
            <ComentarioComponent
              content={comentario.content}
              rating={comentario.rating}
              created_at={comentario.created_at}
              author={comentario.author}
            />
          </ListItem>
        ))}
      </List>
    ) : (
      <Alert>No comments found.</Alert>
    )}
  </section>

  {/* Formulario para añadir comentarios */}
  <section
    style={{
      padding: "30px",
      backgroundColor: "white",
      border: "1px solid #ddd",
      borderRadius: "8px",
      boxShadow: "0 4px 8px rgba(0, 0, 0, 0.1)",
    }}
  >
    <h2 style={{ fontSize: "28px", marginBottom: "20px" }}>Añadir comentario</h2>

    <form id="miFormulario" onSubmit={enviarJSON}>
      <div style={{ marginBottom: "20px" }}>
        <label htmlFor="content" style={{ fontWeight: "bold", fontSize: "18px" }}>
          Contenido:
        </label>
        <input
          type="text"
          id="content"
          name="content"
          required
          style={{
            width: "100%",
            padding: "12px",
            marginTop: "10px",
            border: "1px solid #dd",
            borderRadius: "6px",
            fontSize: "16px",
          }}
        />
      </div>

      <div style={{ marginBottom: "20px" }}>
        <label htmlFor="rating" style={{ fontWeight: "bold", fontSize: "18px" }}>
          Rating:
        </label>
        <input
          type="text"
          id="rating"
          name="rating"
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
        <label htmlFor="author" style={{ fontWeight: "bold", fontSize: "18px" }}>
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

    );
}


export default EntradaPage;
