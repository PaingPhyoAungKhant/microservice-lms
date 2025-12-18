import React, { useState } from 'react';
import type { FormEvent } from 'react';
import { Box, Typography, Link } from '@mui/material';
import { Link as RouterLink, useNavigate } from 'react-router';
import InputField from '../../components/forms/InputField';
import Button from '../../components/common/Button';
import { useAuth } from '../../hooks/useAuth';
import Error from '../../components/common/Error';
import Success from '../../components/common/Success';
import Loading from '../../components/common/Loading';
import { ROUTES } from '../../../shared/constants/routes';

const LoginForm: React.FC = () => {
  const navigate = useNavigate();
  const { login, loading, error } = useAuth();
  const [successMessage, setSuccessMessage] = useState<string | null>(null);
  const [formData, setFormData] = useState({
    email: '',
    password: '',
  });

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData((prev) => ({
      ...prev,
      [name]: value,
    }));
  };

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setSuccessMessage(null);

    try {
      await login(formData.email, formData.password);
      setSuccessMessage('Login successful!');
      navigate(ROUTES.COURSES);
    } catch (err) {
      console.error(err);
    }
  };

  return (
    <Box
      sx={{
        mx: 'auto',
        mt: 4,
        width: '100%',
        maxWidth: 600,
        p: 4,
        border: '1px solid #D9D9D9',
        borderRadius: 2,
        bgcolor: 'background.paper',
        boxShadow: 1,
      }}
    >
      <Typography variant="h4" component="h1" sx={{ mb: 4, textAlign: 'center', fontWeight: 700 }}>
        Login
      </Typography>

      {error && (
        <Error
          message={error.message || 'Failed to login. Please check your credentials.'}
          onRetry={() => {}}
          fullWidth
        />
      )}

      {successMessage && (
        <Success
          message={successMessage}
          autoDismiss={false}
          onDismiss={() => setSuccessMessage(null)}
          fullWidth
        />
      )}

      <form onSubmit={handleSubmit}>
        <InputField
          label="Email"
          type="email"
          name="email"
          placeholder="Enter your email"
          value={formData.email}
          onChange={handleChange}
          required
          disabled={loading}
        />

        <InputField
          label="Password"
          type="password"
          name="password"
          placeholder="Enter your password"
          value={formData.password}
          onChange={handleChange}
          required
          disabled={loading}
        />

        <Box sx={{ mt: 3, display: 'flex', flexDirection: 'column', gap: 2 }}>
          <Button type="submit" variant="fill" color="primary" size="lg" fullWidth disabled={loading}>
            {loading ? 'Logging in...' : 'Log in'}
          </Button>
        </Box>
      </form>

      <Box sx={{ mt: 2, textAlign: 'center' }}>
        <Link component={RouterLink} to={ROUTES.FORGOT_PASSWORD} sx={{ color: 'primary.main', fontWeight: 500 }}>
          Forgot Password?
        </Link>
      </Box>

      <Box
        sx={{
          mt: 4,
          pt: 3,
          borderTop: '1px solid #D9D9D9',
          textAlign: 'center',
          fontSize: '0.875rem',
          color: 'text.secondary',
        }}
      >
        <Typography variant="body2">
          Doesn&#39;t have an account?{' '}
          <Link component={RouterLink} to={ROUTES.REGISTER} sx={{ color: 'primary.main', fontWeight: 500 }}>
            Register
          </Link>
        </Typography>
      </Box>
    </Box>
  );
};

export default LoginForm;

