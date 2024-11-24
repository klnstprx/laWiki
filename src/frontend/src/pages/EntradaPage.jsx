
import { useEffect, useState } from "react";
import List from "@mui/material/List";
import ListItem from "@mui/material/ListItem";
import Alert from "@mui/material/Alert";
import Typography from "@mui/material/Typography";
import { getEntrada } from "../api.js";
import { useSearchParams } from 'react-router-dom';


function EntradaPage() {
  
    const [entrada, setEntrada] = useState({});
    const [error, setError] = useState(null);
    
    const [searchParams] = useSearchParams();

    const id = searchParams.get('id');

    useEffect(() => {
        getEntrada(id)
        .then(setEntrada)
        .catch((err) => setError(err.message));
    }, [id]);


    {/*La URL es de este tipo http://localhost:5173/entrada?id=67311bf03399f3b49ccb8072*/}

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


        ------ Contenido de la wiki (se saca de version) -----------
        <br></br>
        ------ Comentarios (pertenecen a la version no a la entrada) ----------

    </div>
    );
}

export default EntradaPage;
