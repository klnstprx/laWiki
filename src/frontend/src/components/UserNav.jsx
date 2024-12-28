import { useEffect, useState } from "react";
import apiRequest from "../api/Api";
import UserMenu from "./UserMenu";
import LoginButton from "./LoginButton";

const UserNav = () => {
  const [user, setUser] = useState(null);

  useEffect(() => {
    const fetchUserInfo = async () => {
      try {
        const userData = await apiRequest("/auth/me");
        setUser(userData);
      } catch (error) {
        console.error("User is not logged in:", error);
      }
    };
    fetchUserInfo();
  }, []);

  const logout = async () => {
    try {
      await apiRequest("/auth/logout");
      setUser(null);
    } catch (error) {
      console.error("Error logging out:", error);
    }
  };

  return (
    <>
      {user ? <UserMenu user={user} logout={logout} /> : <LoginButton />}
    </>
  );
};

export default UserNav;
