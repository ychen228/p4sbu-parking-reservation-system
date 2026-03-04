import React, { useState, useEffect } from "react";
import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import Sbu_Map from "./components/map";
import { User } from "./api/userApi";
import AuthPage from "./components/authPage";
import Profile from "./components/profile";
import AdminEdit from "./components/adminEdit";
import Tickets from "./components/ticket";
const App: React.FC = () => {
  const [user, set_user] = useState<User | null>(null);

  useEffect(() => {
    const storedUser = localStorage.getItem("user");
    if (storedUser) {
      try {
        const parsedUser = JSON.parse(storedUser);
        if (parsedUser) {
          set_user(parsedUser);
        } else {
          console.error("Invalid user data in localStorage");
        }
      } catch (error) {
        console.error("Error parsing user data from localStorage", error);
      }
    }
  }, []);

  const handleSetUser = (newUser: User | null) => {
    set_user(newUser);
    if (newUser) {
      localStorage.setItem("user", JSON.stringify(newUser));
    } else {
      localStorage.removeItem("user");
    }
  };

  return (
    <div>
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<Navigate to={"/login"} replace />} />
          <Route
            path="/login"
            element={<AuthPage user={user} set_user={handleSetUser} />}
          />
          <Route
            path="/home"
            element={user != null && <Sbu_Map user={user} />}
          />
          <Route path="/profile" element={user && <Profile user={user} set_user={handleSetUser} />} />
          <Route path="/tickets" element={user && <Tickets user={user} />} />

          <Route path="/edit" element={user &&  <AdminEdit user={user} />} />
        </Routes>
      </BrowserRouter>
    </div>
  );
};

// const styles = {
//   app: {
//     minHeight: "100vh",
//     backgroundColor: "#f5f5f5",
//   },
//   container: {
//     maxWidth: "1200px",
//     margin: "0 auto",
//     padding: "0 1rem 2rem",
//   },
// };

export default App;
