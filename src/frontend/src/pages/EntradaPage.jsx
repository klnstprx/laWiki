
import { useEffect, useState } from "react";
import List from "@mui/material/List";
import ListItem from "@mui/material/ListItem";
import Alert from "@mui/material/Alert";
import Typography from "@mui/material/Typography";
import { getEntrada, getAllComentariosByVersion, getVersionById } from "../api.js";
import { useSearchParams } from 'react-router-dom';
import ComentarioComponent from "../components/ComentarioComponent.jsx"
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


    {/*La URL es de este tipo http://localhost:5173/entrada?id=67311bf03399f3b49ccb8072&versionID=67311c0143d96ecd81728a94*/}

    return (
    <div>
        <Typography variant="h3" gutterBottom>
            Datos de la wiki
            {error && <Alert severity="error">{error}</Alert>}
        </Typography>
        
        <Typography variant="h6" gutterBottom>
            Titulo: {entrada.title}
        </Typography>
       
        <Typography variant="h6" gutterBottom>
            Autor: {entrada.author}
        </Typography>

        <Typography variant="h6" gutterBottom>
            Fecha de creación: {entrada.created_at}
        </Typography>
        
        <br></br><br></br>
        ------ Contenido de la wiki (se saca de version) -----------
        <br></br><br></br>


        <VersionComponent 
            content = {version.content} 
            editor = {version.editor} 
            created_at = {version.created_at}
            entry_id = {version.entry_id}
        />

        <br></br>
        ------ Comentarios (pertenecen a la version, no a la entrada) ----------
        
        {comentarios.length > 0
        ? (
          <List>
            {comentarios.map((comentario) => (
              <ListItem key={comentarios.id}>
               <ComentarioComponent
                content = {comentario.content}
                rating = {comentario.rating}
                created_at = {comentario.created_at}
                author = {comentario.author}
               />
              </ListItem>
            ))}
          </List>
        )
        : <Alert>No comments found.</Alert>}

    </div>
    );
}

export default EntradaPage;
