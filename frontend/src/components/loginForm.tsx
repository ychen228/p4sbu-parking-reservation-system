import React, { useState } from 'react';
import userApi, { LoginRequest, User } from '../api/userApi';
import "./css_files/login.css";

interface LoginFormProps {
  onLoginSuccess?: (user: User) => void;
}

const LoginForm: React.FC<LoginFormProps> = ({ onLoginSuccess }) => {
  const [credentials, setCredentials] = useState<LoginRequest>({
    username: '',
    password: '',
  });

  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<boolean>(false);
  const [debugResponse, setDebugResponse] = useState<any>(null);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setCredentials((prev) => ({ ...prev, [name]: value }));
    if (error) setError(null);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!credentials.username.trim() || !credentials.password.trim()) {
      setError('Username and password are required');
      return;
    }

    setLoading(true);
    setError(null);

    try {
      const response = await userApi.login(credentials);
      console.log('API Response:', response);
      setDebugResponse(response);
      setSuccess(true);

      if (onLoginSuccess && response.user) {
        if (!response.user.name) {
          response.user.name = { first: 'User', last: '' };
        }
        onLoginSuccess(response.user);
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Login failed');
      setSuccess(false);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="login-container">
      <h2 className="login-heading">Login</h2>

      {success ? (
        <div className="success-message">
          <p>Login successful!</p>
          <button
            className="login-button success-button"
            onClick={() => {
              setSuccess(false);
              setCredentials({ username: '', password: '' });
            }}
          >
            Login with another account
          </button>

          {debugResponse && (
            <div className="debug-panel">
              <h4>Debug - Response Data:</h4>
              <pre>{JSON.stringify(debugResponse, null, 2)}</pre>
            </div>
          )}
        </div>
      ) : (
        <form onSubmit={handleSubmit}>
          {error && (
            <div className="error-message">
              <p>{error}</p>
            </div>
          )}

          <div className="form-group">
            <label className="form-label" htmlFor="username">Username</label>
            <input
              className="form-input"
              type="text"
              id="username"
              name="username"
              value={credentials.username}
              onChange={handleChange}
              disabled={loading}
              placeholder="Enter your username"
              required
            />
          </div>

          <div className="form-group">
            <label className="form-label" htmlFor="password">Password</label>
            <input
              className="form-input"
              type="password"
              id="password"
              name="password"
              value={credentials.password}
              onChange={handleChange}
              disabled={loading}
              placeholder="Enter your password"
              required
            />
          </div>

          <button
            type="submit"
            className={`login-button ${loading ? 'button-disabled' : ''}`}
            disabled={loading}
          >
            {loading ? 'Logging in...' : 'Login'}
          </button>
        </form>
      )}
    </div>
  );
};

export default LoginForm;
