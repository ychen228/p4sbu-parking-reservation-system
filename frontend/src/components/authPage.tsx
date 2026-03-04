import React, { useState } from 'react';
import LoginForm from './loginForm';
import RegisterForm from './registerForm';
import { User } from '../api/userApi';
import { useNavigate } from "react-router-dom";
import "./css_files/authPage.css";

interface AuthPageProps {
  user: User | null;
  set_user: (user: User | null) => void;
}

const AuthPage: React.FC<AuthPageProps> = ({ user, set_user }) => {
  const navigate = useNavigate();
  const [activeTab, setActiveTab] = useState<'login' | 'register'>('login');

  const handleLoginSuccess = (current_user: User) => {
    set_user(current_user);
    localStorage.setItem("user", JSON.stringify(current_user)); // Save to localStorage
    navigate('/home');
    console.log('Logged in user:', current_user);
  };

  const handleRegisterSuccess = (userId: string) => {
    console.log('Registered new user with ID:', userId);
    setActiveTab('login');
  };


  return (
    <div className="auth-page">
      <h1 className="heading">P4SBU Authentication</h1>

      <div className="tab-navigation">
        <button
          className={`tab-button ${activeTab === 'login' ? 'active' : ''}`}
          onClick={() => setActiveTab('login')}
        >
          Login
        </button>
        <button
          className={`tab-button ${activeTab === 'register' ? 'active' : ''}`}
          onClick={() => setActiveTab('register')}
        >
          Register
        </button>
      </div>

      <div>
        {activeTab === 'login' ? (
          <LoginForm onLoginSuccess={handleLoginSuccess} />
        ) : (
          <RegisterForm onRegisterSuccess={handleRegisterSuccess} />
        )}
      </div>
    </div>
  );
};

export default AuthPage;
