import PropTypes from "prop-types";
import Header from "../components/Header";
import Footer from "../components/Footer";
import ToastMessage from "../components/ToastMessage";
import { useToast } from "../context/ToastContext.1";
import { Outlet } from "react-router-dom";

function MainLayout({ children }) {
  return (
    <div className="d-flex flex-column vh-100">
      <Header />
      <div className="flex-grow-1 pb-5">
        {children}
        <Outlet />
        <ToastMessagesLayout />
      </div>
      <Footer />
    </div>
  );
}

const ToastMessagesLayout = () => {
  const { toast, hideToast } = useToast();

  return (
    <ToastMessage
      show={toast.show}
      onClose={hideToast}
      message={toast.message}
      severity={toast.severity}
    />
  );
};

MainLayout.propTypes = {
  children: PropTypes.node.isRequired,
};

export default MainLayout;
