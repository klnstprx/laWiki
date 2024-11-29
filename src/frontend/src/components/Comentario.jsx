import '../styles/Comentario.css';
import PropTypes from 'prop-types';

const Comentario = ({ id, content, rating, created_at, author }) => {
    
    async function eliminarComentario(event) {
        console.log("Enviando formulario...");
        // Prevenir el envío normal del formulario
        event.preventDefault();
    
        try {
           await fetch("http://localhost:8000/api/comments/id?id=" + id.toString(), {
            method: "DELETE",
            headers: {
              "Content-Type": "application/json",
            },
          });

        } catch (error) {
          console.error("Error al enviar:", error);
        }
      }
    
    
    
    return (
        <div className="comentario">
            <div className="comentario-avatar">
                <img
                    src={`https://ui-avatars.com/api/?name=${author}&background=random`}
                    alt={author}
                    className="comentario-avatar-img"
                />
            </div>
            <div className="comentario-body">
                <div className="comentario-header">
                    <h4 className="comentario-author">{author}</h4>
                    <span className="comentario-date">
                        {new Date(created_at).toLocaleDateString()}
                    </span>
                </div>
                <p className="comentario-content" style={{color:"black"}}>{content}</p>
                <div className="comentario-footer">
                    <span className="comentario-rating">Rating: {rating}/5</span>
                    <button className="comentario-delete" onClick={eliminarComentario}>
                        Borrar
                    </button>
                </div>
            </div>
        </div>
    );
};
Comentario.propTypes = {
    id:PropTypes.string.isRequired,
    content: PropTypes.string.isRequired,
    rating: PropTypes.number.isRequired,
    created_at: PropTypes.string.isRequired,
    author: PropTypes.string.isRequired,
    onDelete: PropTypes.func.isRequired,
};


export default Comentario;