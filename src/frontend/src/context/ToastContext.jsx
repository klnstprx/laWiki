import { useState } from "react";
import PropTypes from "prop-types";
import { ToastContext } from "./ToastContext.1";

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
    </ToastContext.Provider>
  );
};

ToastProvider.propTypes = {
  children: PropTypes.node.isRequired,
};
