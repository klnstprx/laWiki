import { createContext, useContext, useState } from "react";
import PropTypes from "prop-types";

export const ToastContext = createContext();
export const useToast = () => useContext(ToastContext);

export const ToastProvider = ({ children }) => {
  const [toast, setToast] = useState({
    show: false,
    message: "",
    color: "bg-info",
  });

  const showToast = (message, color = "bg-info") => {
    setToast({ show: true, message, color });
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
