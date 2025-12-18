import React, { useState } from 'react';
import type { FormEvent } from 'react';
import { Box, Typography } from '@mui/material';
import InputField from '../../components/forms/InputField';
import Button from '../../components/common/Button';
import Error from '../../components/common/Error';
import Success from '../../components/common/Success';

interface ForgotPasswordFormProps {
  onSubmit: (email: string) => Promise<void>;
  loading?: boolean;
  error?: string | null;
}

const ForgotPasswordForm: React.FC<ForgotPasswordFormProps> = ({ onSubmit, loading = false, error }) => {
  const [email, setEmail] = useState('');
  const [successMessage, setSuccessMessage] = useState<string | null>(null);

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setSuccessMessage(null);

    try {
      await onSubmit(email);
      setSuccessMessage('OTP has been sent to your email. Please check your inbox.');
    } catch (err) {
      console.error(err);
    }
  };

  return (
    <Box>
      <Typography variant="h6" component="h2" sx={{ mb: 2, textAlign: 'center' }}>
        Forgot Password
      </Typography>
      <Typography variant="body2" color="text.secondary" sx={{ mb: 3, textAlign: 'center' }}>
        Enter your email address and we&apos;ll send you an OTP to reset your password.
      </Typography>

      {error && (
        <Error
          message={error || 'Failed to send OTP. Please try again.'}
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
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          required
          disabled={loading}
        />

        <Box sx={{ mt: 3 }}>
          <Button type="submit" variant="fill" color="primary" size="lg" fullWidth disabled={loading}>
            {loading ? 'Sending...' : 'Send OTP'}
          </Button>
        </Box>
      </form>
    </Box>
  );
};

export default ForgotPasswordForm;

