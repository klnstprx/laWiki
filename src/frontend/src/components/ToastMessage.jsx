import { Snackbar, Alert } from "@mui/material";
import PropTypes from "prop-types";

function ToastMessage({
  show,
  message,
  onClose,
  delay = 3000,
  severity = "info",
}) {
  return (
    <Snackbar
      open={show}
      autoHideDuration={delay}
      onClose={onClose}
      anchorOrigin={{ vertical: "top", horizontal: "right" }}
    >
      <Alert onClose={onClose} severity={severity} sx={{ width: "100%" }}>
        {message}
      </Alert>
    </Snackbar>
  );
}

ToastMessage.propTypes = {
  show: PropTypes.bool.isRequired,
  message: PropTypes.string.isRequired,
  onClose: PropTypes.func.isRequired,
  delay: PropTypes.number,
  severity: PropTypes.oneOf(["error", "warning", "info", "success"]),
};

export default ToastMessage;
