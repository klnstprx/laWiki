import { Modal, Button } from "react-bootstrap";
import PropTypes from 'prop-types';

function ConfirmationModal({ show, handleClose, handleConfirm, message }) {
  return (
    <Modal show={show} onHide={handleClose}>
      <Modal.Header closeButton>
        <Modal.Title>Confirmar Acccion</Modal.Title>
      </Modal.Header>
      <Modal.Body>
        {message || "¿Estás seguro de que quieres hacer esto?"}
      </Modal.Body>
      <Modal.Footer>
        <Button variant="primary" onClick={handleConfirm}>
          Confirmar
        </Button>
        <Button variant="light" className="btn btn-primary" onClick={handleClose}>
          Cancelar
        </Button>
      </Modal.Footer>
    </Modal>
  );
}
ConfirmationModal.propTypes = {
  show: PropTypes.bool.isRequired,
  handleClose: PropTypes.func.isRequired,
  handleConfirm: PropTypes.func.isRequired,
  message: PropTypes.string
};

export default ConfirmationModal;
