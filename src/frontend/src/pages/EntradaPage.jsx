
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

        <br></br><br></br><br></br><br></br>


    <Typography variant="h6" gutterBottom>
        Añadir comentario
    </Typography>
    
    
    <form id="miFormulario" onSubmit={enviarJSON}>
        <label htmlFor="content">Contenido:</label>
        <input type="text" id="content" name="content" required/>
        <br></br>
        <label htmlFor="rate">Rating:</label>
        <input type="text" id="rating" name="rating" required/>
        <br></br>
        <label htmlFor="author">Author:</label>
        <input type="text" id="author" name="author" required/>
        <br></br>
       
        <button type="submit">Enviar</button>
    </form>

    </div>
    );
}


export default EntradaPage;
