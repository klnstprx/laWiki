
import { useEffect, useState } from "react";
import List from "@mui/material/List";
import ListItem from "@mui/material/ListItem";
import Alert from "@mui/material/Alert";
import Typography from "@mui/material/Typography";
import { getEntrada } from "../api.js";
import { useParams } from 'react-router-dom';


function EntradaPage() {
  
    const [entrada, setEntrada] = useState({});
    const [error, setError] = useState(null);
    
    const { id } = useParams();

    useEffect(() => {
        getEntrada(id)
        .then(setEntrada)
        .catch((err) => setError(err.message));
    }, [id]);


    return (
    <div>
        <Typography variant="h1" gutterBottom>
            {entrada.title}
        </Typography>

    </div>
    );
}

export default EntradaPage;
