import React from 'react';
import { Box } from '@mui/material';
import LoginForm from '../features/auth/LoginForm';

export default function Login() {
  return (
    <Box
      sx={{
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        px: { xs: 2, md: 8 },
        py: 4,
      }}
    >
      <LoginForm />
    </Box>
  );
}

