import React from 'react';
import { Box } from '@mui/material';
import DashboardLoginForm from '../../features/dashboard/auth/LoginForm';

export default function DashboardLogin() {
  return (
    <Box
      sx={{
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        minHeight: '100vh',
        bgcolor: 'background.default',
      }}
    >
      <DashboardLoginForm />
    </Box>
  );
}

