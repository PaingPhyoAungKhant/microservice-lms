import React from 'react';
import { Box } from '@mui/material';
import RegisterForm from '../features/auth/RegisterForm';

export default function Register() {
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
      <RegisterForm />
    </Box>
  );
}

