import React, { useState } from 'react';
import type { FormEvent } from 'react';
import { Box, Typography } from '@mui/material';
import InputField from '../../../components/forms/InputField';
import Button from '../../../components/common/Button';
import { useDashboardAuth } from '../../../hooks/useDashboardAuth';
import Error from '../../../components/common/Error';
import Success from '../../../components/common/Success';

const DashboardLoginForm: React.FC = () => {
  const { login, loading, error } = useDashboardAuth();
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
      setSuccessMessage('Login successful! Redirecting...');
    } catch (err) {
      console.error(err);
    }
  };

  return (
    <Box
      sx={{
        mx: 'auto',
        mt: 8,
        width: '100%',
        maxWidth: 500,
        p: 4,
        border: '1px solid #D9D9D9',
        borderRadius: 2,
        bgcolor: 'background.paper',
        boxShadow: 3,
      }}
    >
      <Typography variant="h4" component="h1" sx={{ mb: 4, textAlign: 'center', fontWeight: 700 }}>
        Dashboard Login
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

        <Box sx={{ mt: 3 }}>
          <Button type="submit" variant="fill" color="primary" size="lg" fullWidth disabled={loading}>
            {loading ? 'Logging in...' : 'Log in'}
          </Button>
        </Box>
      </form>
    </Box>
  );
};

export default DashboardLoginForm;

