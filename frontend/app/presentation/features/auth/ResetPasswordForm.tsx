import React, { useState } from 'react';
import type { FormEvent } from 'react';
import { Box, Typography } from '@mui/material';
import InputField from '../../components/forms/InputField';
import Button from '../../components/common/Button';
import Error from '../../components/common/Error';
import Success from '../../components/common/Success';

interface ResetPasswordFormProps {
  token: string;
  onSubmit: (token: string, newPassword: string) => Promise<void>;
  onBack: () => void;
  loading?: boolean;
  error?: string | null;
}

const ResetPasswordForm: React.FC<ResetPasswordFormProps> = ({
  token,
  onSubmit,
  onBack,
  loading = false,
  error,
}) => {
  const [formData, setFormData] = useState({
    newPassword: '',
    confirmPassword: '',
  });
  const [successMessage, setSuccessMessage] = useState<string | null>(null);
  const [passwordError, setPasswordError] = useState<string | null>(null);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData((prev) => ({
      ...prev,
      [name]: value,
    }));
    setPasswordError(null);
  };

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setPasswordError(null);
    setSuccessMessage(null);

    if (formData.newPassword !== formData.confirmPassword) {
      setPasswordError('Passwords do not match');
      return;
    }

    if (formData.newPassword.length < 8) {
      setPasswordError('Password must be at least 8 characters');
      return;
    }

    try {
      await onSubmit(token, formData.newPassword);
      setSuccessMessage('Password reset successfully! Redirecting to login...');
    } catch (err) {
      console.error(err);
    }
  };

  return (
    <Box>
      <Typography variant="h6" component="h2" sx={{ mb: 2, textAlign: 'center' }}>
        Reset Password
      </Typography>
      <Typography variant="body2" color="text.secondary" sx={{ mb: 3, textAlign: 'center' }}>
        Enter your new password
      </Typography>

      {error && (
        <Error
          message={error || 'Failed to reset password. Please try again.'}
          onRetry={() => {}}
          fullWidth
        />
      )}

      {passwordError && (
        <Error
          message={passwordError}
          onRetry={() => setPasswordError(null)}
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
          label="New Password"
          type="password"
          name="newPassword"
          placeholder="Enter new password"
          value={formData.newPassword}
          onChange={handleChange}
          required
          disabled={loading}
          helperText="Password must be at least 8 characters"
        />

        <InputField
          label="Confirm Password"
          type="password"
          name="confirmPassword"
          placeholder="Confirm new password"
          value={formData.confirmPassword}
          onChange={handleChange}
          required
          disabled={loading}
        />

        <Box sx={{ mt: 3, display: 'flex', gap: 2 }}>
          <Button
            type="button"
            variant="outline"
            onClick={onBack}
            disabled={loading}
            fullWidth
          >
            Back
          </Button>
          <Button
            type="submit"
            variant="fill"
            color="primary"
            size="lg"
            fullWidth
            disabled={loading || !formData.newPassword || !formData.confirmPassword}
          >
            {loading ? 'Resetting...' : 'Reset Password'}
          </Button>
        </Box>
      </form>
    </Box>
  );
};

export default ResetPasswordForm;

