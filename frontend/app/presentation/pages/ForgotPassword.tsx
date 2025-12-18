import React, { useState } from 'react';
import { Box, Typography, Link } from '@mui/material';
import { Link as RouterLink, useNavigate } from 'react-router';
import {
  useForgotPasswordMutation,
  useVerifyOTPMutation,
  useResetPasswordMutation,
} from '../../infrastructure/api/rtk/authApi';
import ForgotPasswordForm from '../features/auth/ForgotPasswordForm';
import VerifyOTPForm from '../features/auth/VerifyOTPForm';
import ResetPasswordForm from '../features/auth/ResetPasswordForm';
import Error from '../components/common/Error';
import { getErrorMessage } from '../../infrastructure/api/utils';
import { ROUTES } from '../../shared/constants/routes';

type Step = 'email' | 'otp' | 'reset';

const ForgotPassword: React.FC = () => {
  const navigate = useNavigate();
  const [step, setStep] = useState<Step>('email');
  const [email, setEmail] = useState('');
  const [resetToken, setResetToken] = useState<string | null>(null);

  const [forgotPassword, { isLoading: forgotPasswordLoading, error: forgotPasswordError }] =
    useForgotPasswordMutation();
  const [verifyOTP, { isLoading: verifyOTPLoading, error: verifyOTPError, reset: resetVerifyOTP }] = useVerifyOTPMutation();
  const [resetPassword, { isLoading: resetPasswordLoading, error: resetPasswordError }] =
    useResetPasswordMutation();

  const handleForgotPassword = async (emailValue: string) => {
    try {
      await forgotPassword({ email: emailValue }).unwrap();
      setEmail(emailValue);
      setStep('otp');
    } catch (err) {
      throw err;
    }
  };

  const handleVerifyOTP = async (emailValue: string, otp: string): Promise<string> => {
    resetVerifyOTP();
    try {
      const result = await verifyOTP({ email: emailValue, otp }).unwrap();
      if (result.isValid && result.passwordResetToken) {
        setResetToken(result.passwordResetToken);
        setStep('reset');
        return result.passwordResetToken;
      } else {
        throw new Error(result.errorMessage || 'Invalid OTP. Please try again.');
      }
    } catch (err) {
      throw err;
    }
  };

  const handleResetPassword = async (tokenValue: string, newPassword: string) => {
    try {
      await resetPassword({ token: tokenValue, new_password: newPassword }).unwrap();
      setTimeout(() => {
        navigate(ROUTES.LOGIN);
      }, 2000);
    } catch (err) {
      throw err;
    }
  };

  const handleBack = () => {
    if (step === 'otp') {
      setStep('email');
      setEmail('');
      resetVerifyOTP();
    } else if (step === 'reset') {
      setStep('otp');
      setResetToken(null);
      resetVerifyOTP();
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
      {step === 'email' && (
        <ForgotPasswordForm
          onSubmit={handleForgotPassword}
          loading={forgotPasswordLoading}
          error={forgotPasswordError ? getErrorMessage(forgotPasswordError, 'Failed to send OTP. Please try again.') : null}
        />
      )}

      {step === 'otp' && (
        <VerifyOTPForm
          email={email}
          onSubmit={handleVerifyOTP}
          onBack={handleBack}
          loading={verifyOTPLoading}
          error={verifyOTPError ? getErrorMessage(verifyOTPError, 'Invalid OTP. Please try again.') : null}
        />
      )}

      {step === 'reset' && resetToken && (
        <ResetPasswordForm
          token={resetToken}
          onSubmit={handleResetPassword}
          onBack={handleBack}
          loading={resetPasswordLoading}
          error={resetPasswordError ? getErrorMessage(resetPasswordError, 'Failed to reset password. Please try again.') : null}
        />
      )}

      <Box sx={{ mt: 3, textAlign: 'center' }}>
        <Link component={RouterLink} to={ROUTES.LOGIN} sx={{ color: 'primary.main', fontWeight: 500 }}>
          Back to Login
        </Link>
      </Box>
    </Box>
  );
};

export default ForgotPassword;

