import { useEffect, useState } from "react";
import { useSearchParams } from "react-router-dom";
import MainLayout from "../layout/MainLayout.jsx";
// import { useToast } from "../context/ToastContext.1.jsx";
import List from "@mui/material/List";
import ListItem from "@mui/material/ListItem";
import Alert from "@mui/material/Alert";
import VersionCard from "../components/VersionCard.jsx";
import { searchVersions } from "../api/VersionApi.js";

function VersionPage() {
  const [versiones, setVersiones] = useState([]);
  const [error, setError] = useState(null);

  const [searchParams] = useSearchParams();
  const entryID = searchParams.get("entryID");

  useEffect(() => {
    searchVersions({ entryID: entryID })
      .then(setVersiones)
      .catch((err) => setError(err.message));
  }, [entryID]);

  {
    /* http://localhost:5173/versiones?entryID=67311bf03399f3b49ccb8072 */
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
          height: "80vh",
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
          {!error && versiones.length > 0 ? (
            <List>
              {versiones.map((version) => (
                <ListItem
                  key={version.id}
                  style={{
                    borderBottom: "1px solid #ddd",
                    padding: "15px 0",
                  }}
                >
                  <VersionCard
                    entradaID={entryID}
                    versionId={version.id}
                    editor={version.editor}
                    created_at={version.created_at}
                  />
                </ListItem>
              ))}
            </List>
          ) : (
            <Alert>No versions found.</Alert>
          )}
        </section>
      </div>
    </MainLayout>
  );
}

export default VersionPage;
