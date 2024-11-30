

const VersionCard = ({ entradaID, versionId, editor, created_at }) => {
    
  
    return (
      <div>
        <span>Fecha: {created_at}&nbsp;&nbsp;&nbsp;&nbsp; Editor: {editor} &nbsp;&nbsp; <a href={`http://localhost:5173/entrada?id=${entradaID}&versionID=${versionId}`}>Ver</a></span>
      </div>
    );
  };


export default VersionCard;
