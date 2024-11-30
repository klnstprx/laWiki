import { useEffect, useState } from "react";
import { useRef } from "react";
import List from "@mui/material/List";
import ListItem from "@mui/material/ListItem";
import Alert from "@mui/material/Alert";
import Typography from "@mui/material/Typography";
import { searchComments } from "../api/CommentApi.js";
import { getEntry } from "../api/EntryApi.js";
import { getVersion } from "../api/VersionApi.js";
import { postComment } from "../api/CommentApi.js";
import { deleteComment } from "../api/CommentApi.js";
import { useSearchParams } from "react-router-dom";
import Comentario from "../components/Comentario.jsx";
import Version from "../components/Version.jsx";
import MainLayout from "../layout/MainLayout.jsx";
import ConfirmationModal from "../components/ConfirmationModal.jsx";
//import { useToast } from "../context/ToastContext.jsx";

function EntradaPage() {
  const [entrada, setEntrada] = useState({});
  const [version, setVersion] = useState({});
  const [comentarios, setComentarios] = useState([]);
  const [entryError, setEntryError] = useState(null);
  const [commentsError, setCommentsError] = useState(null);
  const [versionError, setVersionError] = useState(null);

  const [searchParams] = useSearchParams();
  const id = searchParams.get("id");
  const versionID = searchParams.get("versionID");

  const [showModal, setShowModal] = useState(false);
  const [pendingComment, setPendingComment] = useState(null);
  //const { showToast } = useToast();
  const formRef = useRef(null);

  const handleClose = () => {
    setShowModal(false);
    setPendingComment(null);
    //showToast("El comentario no se ha creado", "bg-danger");
  };

  const handleConfirm = async () => {
    setShowModal(false);
    try {
      const result = await postComment(pendingComment);

      setComentarios((prevComentarios) => [...prevComentarios, result]);

      formRef.current.reset();
      setPendingComment(null);

      //showToast("El comentario se ha creado correctamente!", "bg-success");
    } catch (error) {
      console.error("Error al enviar:", error);
      //showToast("Error al enviar el comentario", "bg-danger");
    }
  };

  const handleDeleteComment = async (commentId) => {
    try {
      await deleteComment(commentId);
      // Update state to remove the deleted comment
      setComentarios((prevComentarios) =>
        prevComentarios.filter((comment) => comment.id !== commentId),
      );
      //showToast("Comentario eliminado correctamente", "bg-success");
    } catch (error) {
      console.error("Error al eliminar el comentario:", error);
      //showToast("Error al eliminar el comentario", "bg-danger");
    }
  };

  useEffect(() => {
    getEntry(id)
      .then(setEntrada)
      .catch((err) => setEntryError(err.message));
  }, [id]);

  useEffect(() => {
    searchComments({ versionID: versionID })
      .then(setComentarios)
      .catch((err) => setCommentsError(err.message));
  }, [versionID]);

  useEffect(() => {
    getVersion(versionID)
      .then(setVersion)
      .catch((err) => setVersionError(err.message));
  }, [versionID]);

  async function subirComentario(event) {
    event.preventDefault();

    const formData = new FormData(event.target);
    const jsonData = Object.fromEntries(formData.entries());

    jsonData["version_id"] = version.id;
    jsonData["rating"] = parseInt(jsonData["rating"], 10);

    // Store the comment data and show the modal
    setPendingComment(jsonData);
    setShowModal(true);
  }

  {
    /*La URL es de este tipo http://localhost:5173/entrada?id=67311bf03399f3b49ccb8072&versionID=67311c0143d96ecd81728a94 */
  }

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
          width: "94vw",
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
          <h1 style={{ fontSize: "36px", margin: "0" }}>Datos de la Entrada</h1>
          {entryError && (
            <div
              style={{
                backgroundColor: "#e57373",
                padding: "15px",
                marginTop: "15px",
                borderRadius: "4px",
              }}
            >
              <p>Error al cargar la entrada: {entryError}</p>
            </div>
          )}
        </header>

        {/* Información de la entrada */}
        <section
          style={{
            padding: "30px",
            backgroundColor: "white",
            margin: "20px 0",
            borderRadius: "8px",
          }}
        >
          {versionError && (
            <div
              style={{
                backgroundColor: "#e57373",
                padding: "15px",
                marginBottom: "15px",
                borderRadius: "4px",
              }}
            >
              <p>Error al cargar la versión: {versionError}</p>
            </div>
          )}
          {!entryError && (
            <>
              <Typography
                variant="h6"
                style={{ fontWeight: "bold", marginBottom: "10px" }}
              >
                Título:{" "}
                <span style={{ fontWeight: "normal" }}>{entrada.title}</span>
              </Typography>
              <Typography
                variant="h6"
                style={{ fontWeight: "bold", marginBottom: "10px" }}
              >
                Autor:{" "}
                <span style={{ fontWeight: "normal" }}>{entrada.author}</span>
              </Typography>
              <Typography
                variant="h6"
                style={{ fontWeight: "bold", marginBottom: "10px" }}
              >
                Fecha de creación:{" "}
                <span style={{ fontWeight: "normal" }}>
                  {entrada.created_at}
                </span>
              </Typography>
              <span><a href={`http://localhost:5173/versiones?entryID=${entrada.id}`}>Ver historial</a></span>
            </>
          )}
        </section>

        {/* Contenido de la versión */}
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
            Contenido de la versión
          </h2>
          {!versionError && (
            <Version
              content={version.content}
              editor={version.editor}
              created_at={version.created_at}
              entry_id={version.entry_id}
            />
          )}
        </section>

        {/* Comentarios */}
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
            Comentarios
          </h2>
          {commentsError && (
            <div
              style={{
                backgroundColor: "#e57373",
                padding: "15px",
                marginBottom: "15px",
                borderRadius: "4px",
              }}
            >
              <p>Error al cargar los comentarios: {commentsError}</p>
            </div>
          )}
          {!commentsError && comentarios.length > 0 ? (
            <List>
              {comentarios.map((comentario) => (
                <ListItem
                  key={comentario.id}
                  style={{
                    borderBottom: "1px solid #ddd",
                    padding: "15px 0",
                  }}
                >
                  <Comentario
                    id={comentario.id}
                    content={comentario.content}
                    rating={comentario.rating}
                    created_at={comentario.created_at}
                    author={comentario.author}
                    onDelete={handleDeleteComment}
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
          <h2 style={{ fontSize: "28px", marginBottom: "20px" }}>
            Añadir comentario
          </h2>

          <form id="miFormulario" ref={formRef} onSubmit={subirComentario}>
            <div style={{ marginBottom: "20px" }}>
              <label
                htmlFor="content"
                style={{ fontWeight: "bold", fontSize: "18px" }}
              >
                Contenido:
                <br />
              </label>
              <textarea
                //type="text"
                id="content"
                name="content"
                required
                style={{
                  width: "70%",
                  padding: "12px",
                  marginTop: "10px",
                  border: "1px solid #dd",
                  borderRadius: "6px",
                  fontSize: "16px",
                }}
              />
            </div>

            <div style={{ marginBottom: "20px" }}>
              <label
                htmlFor="rating"
                style={{ fontWeight: "bold", fontSize: "18px" }}
              >
                Rating:
              </label>
              <input
                type="text"
                id="rating"
                name="rating"
                required
                style={{
                  width: "5%",
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
                  width: "15%",
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
                width: "5%",
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
        message="¿Estás seguro de que quieres crear este comentario?"
        show={showModal}
        handleClose={handleClose}
        handleConfirm={handleConfirm}
      />
    </MainLayout>
  );
}

export default EntradaPage;
