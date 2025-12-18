import React from 'react';
import { Alert, AlertTitle, Box, Button } from '@mui/material';
import ErrorOutlineIcon from '@mui/icons-material/ErrorOutline';

type ErrorProps = {
  message?: string;
  title?: string;
  onRetry?: () => void;
  fullWidth?: boolean;
  severity?: 'error' | 'warning';
};

export default function Error({
  message = 'An error occurred. Please try again.',
  title = 'Error',
  onRetry,
  fullWidth = false,
  severity = 'error',
}: ErrorProps) {
  return (
    <Box sx={{ width: fullWidth ? '100%' : 'auto', p: 2 }}>
      <Alert
        severity={severity}
        icon={<ErrorOutlineIcon />}
        action={
          onRetry ? (
            <Button color="inherit" size="small" onClick={onRetry}>
              Retry
            </Button>
          ) : undefined
        }
        sx={{
          '& .MuiAlert-message': {
            width: '100%',
          },
        }}
      >
        <AlertTitle>{title}</AlertTitle>
        {message}
      </Alert>
    </Box>
  );
}

