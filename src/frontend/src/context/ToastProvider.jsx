import { useState } from "react";
import PropTypes from "prop-types";
import { ToastContext } from "./ToastContext";
import ToastMessage from "../components/ToastMessage";

export const ToastProvider = ({ children }) => {
  const [toast, setToast] = useState({
    show: false,
    message: "",
    severity: "info",
  });

  const showToast = (message, severity = "info") => {
    setToast({ show: true, message, severity });
  };

  const hideToast = () => {
    setToast({ ...toast, show: false });
  };

  return (
    <ToastContext.Provider value={{ toast, showToast, hideToast }}>
      {children}
      <ToastMessage
        show={toast.show}
        message={toast.message}
        onClose={hideToast}
        severity={toast.severity}
      />
    </ToastContext.Provider>
  );
};

ToastProvider.propTypes = {
  children: PropTypes.node.isRequired,
};
