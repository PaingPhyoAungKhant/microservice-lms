import React, { useState } from 'react';
import type { FormEvent } from 'react';
import { Box, Typography } from '@mui/material';
import InputField from '../../components/forms/InputField';
import Button from '../../components/common/Button';
import Error from '../../components/common/Error';
import Success from '../../components/common/Success';

interface VerifyOTPFormProps {
  email: string;
  onSubmit: (email: string, otp: string) => Promise<string>;
  onBack: () => void;
  loading?: boolean;
  error?: string | null;
}

const VerifyOTPForm: React.FC<VerifyOTPFormProps> = ({
  email,
  onSubmit,
  onBack,
  loading = false,
  error,
}) => {
  const [otp, setOtp] = useState('');
  const [successMessage, setSuccessMessage] = useState<string | null>(null);

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setSuccessMessage(null);

    try {
      const token = await onSubmit(email, otp);
      setSuccessMessage('OTP verified successfully! Redirecting...');
    } catch (err) {
      console.error(err);
    }
  };

  return (
    <Box>
      <Typography variant="h6" component="h2" sx={{ mb: 2, textAlign: 'center' }}>
        Verify OTP
      </Typography>
      <Typography variant="body2" color="text.secondary" sx={{ mb: 3, textAlign: 'center' }}>
        Enter the 6-digit OTP sent to {email}
      </Typography>

      {error && (
        <Error
          message={error || 'Invalid OTP. Please try again.'}
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
          label="OTP"
          type="text"
          name="otp"
          placeholder="Enter 6-digit OTP"
          value={otp}
          onChange={(e) => {
            const value = e.target.value.replace(/\D/g, '').slice(0, 6);
            setOtp(value);
          }}
          required
          disabled={loading}
          inputProps={{ maxLength: 6, pattern: '[0-9]{6}' }}
          helperText="Enter the 6-digit code from your email"
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
            disabled={loading || otp.length !== 6}
          >
            {loading ? 'Verifying...' : 'Verify OTP'}
          </Button>
        </Box>
      </form>
    </Box>
  );
};

export default VerifyOTPForm;
  
