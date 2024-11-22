

const ComentarioComponent = ({ content, rating, created_at, author }) => {
	return (
		<div className="comentario">
			<h4 className="comentario-author">{author}</h4>
			<p className="comentario-content">{content}</p>
			<div className="comentario-meta">
				<span className="comentario-rating">Rating: {rating}</span>
				<span className="comentario-date">{new Date(created_at).toLocaleDateString()}</span>
			</div>
    	</div>
	);
};

export default ComentarioComponent;