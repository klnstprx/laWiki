import { createContext, useContext, useEffect, useState } from "react";
import apiRequest from "../api/Api";

const AuthContext = createContext();

export const AuthProvider = ({ children }) => {
  const [user, setUser] = useState(null);

  // Fetch user info on mount
  useEffect(() => {
    const fetchUserInfo = async () => {
      try {
        const userData = await apiRequest("/auth/me");
        setUser(userData);
      } catch (error) {
        // User is not authenticated
        setUser(null);
      }
    };
    fetchUserInfo();
  }, []);

  const logout = async () => {
    try {
      await apiRequest("/auth/logout");
    } catch (error) {
      console.error("Error during logout:", error);
    }
    setUser(null);
  };

  return (
    <AuthContext.Provider value={{ user, setUser, logout }}>
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => useContext(AuthContext);
