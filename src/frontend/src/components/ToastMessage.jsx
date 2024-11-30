import { Toast } from "react-bootstrap";
import PropTypes from "prop-types";

function ToastMessage({
  show,
  message,
  onClose,
  delay = 3000,
  color = "bg-info",
}) {
  return (
    <Toast
      show={show}
      onClose={onClose}
      delay={delay}
      autohide
      className={`${color}`}
      style={{
        position: "fixed",
        top: "20px",
        right: "20px",
        zIndex: 9999,
      }}
    >
      <Toast.Header>
        <strong className="m-auto">Notificacion</strong>
      </Toast.Header>
      <Toast.Body style={{ color: "black" }}>{message}</Toast.Body>
    </Toast>
  );
}

ToastMessage.propTypes = {
  show: PropTypes.bool.isRequired,
  message: PropTypes.string.isRequired,
  onClose: PropTypes.func.isRequired,
  delay: PropTypes.number,
  color: PropTypes.string,
};

export default ToastMessage;
