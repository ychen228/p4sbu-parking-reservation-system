import React, { useState } from 'react';
import userApi, { RegisterRequest } from '../api/userApi';
import "./css_files/register.css";

interface RegisterFormProps {
  onRegisterSuccess?: (userId: string) => void;
}

const RegisterForm: React.FC<RegisterFormProps> = ({ onRegisterSuccess }) => {
  const [formData, setFormData] = useState<RegisterRequest>({
    username: '',
    passwordHash: '',
    name: {
      first: '',
      last: '',
    },
    sbuId: '',
  });

  const [confirmPassword, setConfirmPassword] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState(false);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;

    if (name === 'firstName') {
      setFormData(prev => ({ ...prev, name: { ...prev.name, first: value } }));
    } else if (name === 'lastName') {
      setFormData(prev => ({ ...prev, name: { ...prev.name, last: value } }));
    } else if (name === 'confirmPassword') {
      setConfirmPassword(value);
    } else {
      setFormData(prev => ({ ...prev, [name]: value }));
    }

    if (error) setError(null);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!formData.username.trim() || !formData.passwordHash.trim()) {
      setError('All fields are required');
      return;
    }

    if (formData.passwordHash !== confirmPassword) {
      setError('Passwords do not match');
      return;
    }

    setLoading(true);
    setError(null);

    try {
      const response = await userApi.register(formData);
      setSuccess(true);
      onRegisterSuccess?.(response.id);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Registration failed');
      setSuccess(false);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="form-container">
      <h2 className="form-heading">Create an Account</h2>

      {success ? (
        <div className="success-message">
          <p>Registration successful! You can now log in with your credentials.</p>
          <button
            className="form-button success-button"
            onClick={() => {
              setSuccess(false);
              setFormData({
                username: '',
                passwordHash: '',
                name: { first: '', last: '' },
                sbuId: '',
              });
              setConfirmPassword('');
            }}
          >
            Register another account
          </button>
        </div>
      ) : (
        <form onSubmit={handleSubmit}>
          {error && <div className="error-message"><p>{error}</p></div>}

          <div className="form-row">
            <div className="form-group">
              <label htmlFor="firstName" className="form-label">First Name</label>
              <input
                type="text"
                id="firstName"
                name="firstName"
                value={formData.name.first}
                onChange={handleChange}
                disabled={loading}
                required
                className="form-input"
              />
            </div>

            <div className="form-group">
              <label htmlFor="lastName" className="form-label">Last Name</label>
              <input
                type="text"
                id="lastName"
                name="lastName"
                value={formData.name.last}
                onChange={handleChange}
                disabled={loading}
                required
                className="form-input"
              />
            </div>
          </div>

          <div className="form-group">
            <label htmlFor="sbuId" className="form-label">SBU ID</label>
            <input
              type="text"
              id="sbuId"
              name="sbuId"
              value={formData.sbuId}
              onChange={handleChange}
              disabled={loading}
              required
              className="form-input"
            />
          </div>

          <div className="form-group">
            <label htmlFor="username" className="form-label">Username</label>
            <input
              type="text"
              id="username"
              name="username"
              value={formData.username}
              onChange={handleChange}
              disabled={loading}
              required
              className="form-input"
            />
          </div>

          <div className="form-group">
            <label htmlFor="passwordHash" className="form-label">Password</label>
            <input
              type="password"
              id="passwordHash"
              name="passwordHash"
              value={formData.passwordHash}
              onChange={handleChange}
              disabled={loading}
              required
              className="form-input"
            />
          </div>

          <div className="form-group">
            <label htmlFor="confirmPassword" className="form-label">Confirm Password</label>
            <input
              type="password"
              id="confirmPassword"
              name="confirmPassword"
              value={confirmPassword}
              onChange={handleChange}
              disabled={loading}
              required
              className="form-input"
            />
          </div>

          <button
            type="submit"
            disabled={loading}
            className={`form-button ${loading ? 'form-button-disabled' : ''}`}
          >
            {loading ? 'Registering...' : 'Register'}
          </button>
        </form>
      )}
    </div>
  );
};

export default RegisterForm;
