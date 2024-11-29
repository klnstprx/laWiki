import '../styles/Comentario.css';
import PropTypes from 'prop-types';

const Comentario = ({ content, rating, created_at, author, onDelete }) => {
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
                <p className="comentario-content">{content}</p>
                <div className="comentario-footer">
                    <span className="comentario-rating">Rating: {rating}/5</span>
                    <button className="comentario-delete" onClick={onDelete}>
                        Borrar
                    </button>
                </div>
            </div>
        </div>
    );
};
Comentario.propTypes = {
    content: PropTypes.string.isRequired,
    rating: PropTypes.number.isRequired,
    created_at: PropTypes.string.isRequired,
    author: PropTypes.string.isRequired,
    onDelete: PropTypes.func.isRequired,
};

/*
const Comentario = ({ content, rating, created_at, author }) => {
	return (
		<div className="comentario">
			<h4 className="comentario-author">{author}</h4>
			<p className="comentario-content">{content}</p>
			<div className="comentario-meta">
				<p className="comentario-rating">Rating: {rating}</p>
				<span className="comentario-date">Fecha: {new Date(created_at).toLocaleDateString()}</span>
			</div>
    	</div>
	);
};

*/

export default Comentario;