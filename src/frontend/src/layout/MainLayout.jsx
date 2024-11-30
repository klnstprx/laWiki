import Header from "../components/Header";
import Footer from "../components/Footer";
import ToastMessage from "../components/ToastMessage";
import { useToast } from "../context/ToastContext.1";
import { Outlet } from "react-router-dom";

function MainLayout() {
  return (
    <div className="d-flex flex-column vh-100">
      <Header />
      <div className="flex-grow-1 pb-5">
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

export default MainLayout;
