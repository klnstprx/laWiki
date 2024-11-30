import { Link } from "react-router-dom";
import PropTypes from "prop-types";

const EntradaCard = ({ id, title, author, createdAt, onEntradaClick }) => {
  const handleClick = () => {
    if (onEntradaClick) {
      onEntradaClick(id);
    }
  };

  return (
    <Link
      to={`/entrada?id=${id}`}
      onClick={handleClick}
      className="custom-link m-0 p-0"
    >
      <div className="border border-dark-subtle bg-light text-dark mb-2 mx-0 p-0">
        <div className="custom-link-container m-0 p-0">
          <div className="p-3 d-flex align-items-center">
            <span className="h4">{title}</span>
          </div>
          <div className="p-3">
            <div className="row">
              <div className="col m-auto">
                <span className="h5">Author</span>
                <p>{author}</p>
              </div>
              <div className="col m-auto">
                <span className="h5">Created At</span>
                <p>{createdAt}</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </Link>
  );
};

EntradaCard.propTypes = {
  id: PropTypes.string.isRequired,
  title: PropTypes.string.isRequired,
  author: PropTypes.string.isRequired,
  createdAt: PropTypes.string.isRequired,
  onEntradaClick: PropTypes.func,
};

export default EntradaCard;
