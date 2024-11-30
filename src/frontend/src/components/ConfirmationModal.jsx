import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogContentText,
  DialogActions,
  Button,
} from "@mui/material";
import PropTypes from "prop-types";

function ConfirmationModal({ show, handleClose, handleConfirm, message }) {
  return (
    <Dialog open={show} onClose={handleClose}>
      <DialogTitle>Confirmar Acción</DialogTitle>
      <DialogContent>
        <DialogContentText>
          {message || "¿Estás seguro de que quieres hacer esto?"}
        </DialogContentText>
      </DialogContent>
      <DialogActions>
        <Button onClick={handleClose}>Cancelar</Button>
        <Button onClick={handleConfirm} color="primary" variant="contained">
          Confirmar
        </Button>
      </DialogActions>
    </Dialog>
  );
}

ConfirmationModal.propTypes = {
  show: PropTypes.bool.isRequired,
  handleClose: PropTypes.func.isRequired,
  handleConfirm: PropTypes.func.isRequired,
  message: PropTypes.string,
};

export default ConfirmationModal;
