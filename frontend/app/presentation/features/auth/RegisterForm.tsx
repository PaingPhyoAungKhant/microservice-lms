import React, { useState } from 'react';
import type { FormEvent } from 'react';
import { Box, Typography, Link } from '@mui/material';
import { Link as RouterLink, useNavigate } from 'react-router';
import InputField from '../../components/forms/InputField';
import Button from '../../components/common/Button';
import { useAuth } from '../../hooks/useAuth';
import Error from '../../components/common/Error';
import Success from '../../components/common/Success';
import { ROUTES } from '../../../shared/constants/routes';

const RegisterForm: React.FC = () => {
  const navigate = useNavigate();
  const { register, loading, error } = useAuth();
  const [successMessage, setSuccessMessage] = useState<string | null>(null);
  const [formData, setFormData] = useState({
    username: '',
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
      await register(formData.username, formData.email, formData.password);
      setSuccessMessage('Registration successful! Please check your email to verify your account. Redirecting to login...');
      setTimeout(() => {
        navigate(ROUTES.LOGIN);
      }, 2000);
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
        Register Now
      </Typography>

      {error && (
        <Error
          message={error.message || 'Failed to register. Please try again.'}
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
          label="Username"
          name="username"
          placeholder="Enter your username"
          value={formData.username}
          onChange={handleChange}
          required
          disabled={loading}
          helperText="Username must be at least 3 characters"
        />

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
          helperText="Password must be at least 8 characters"
        />

        <Box sx={{ mt: 3 }}>
          <Button type="submit" variant="fill" color="primary" size="lg" fullWidth disabled={loading}>
            {loading ? 'Registering...' : 'Register'}
          </Button>
        </Box>
      </form>

      <Box sx={{ mt: 2, textAlign: 'center' }}>
        <Typography variant="body2">
          Already has an account?{' '}
          <Link component={RouterLink} to={ROUTES.LOGIN} sx={{ color: 'primary.main', fontWeight: 500 }}>
            Log in
          </Link>
        </Typography>
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
        <Typography variant="body2">By registering, you agree to our Terms of Service</Typography>
      </Box>
    </Box>
  );
};

export default RegisterForm;

