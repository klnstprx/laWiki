import React from "react";
import Typography from "@mui/material/Typography";

const VersionComponent = (props) => {
  
    const { content, editor, created_at, entry_id  } = props;

    return (
    <div>
        <Typography variant="h6" gutterBottom>
            contenido: {content}
        </Typography>
       
        <Typography variant="h6" gutterBottom>
            Editor: {editor}
        </Typography>

        <Typography variant="h6" gutterBottom>
            Fecha de creación: {created_at}
        </Typography>

    </div>
    );
};

export default VersionComponent;






