import React, { useEffect } from 'react';
import { Alert, AlertTitle, Box } from '@mui/material';
import CheckCircleOutlineIcon from '@mui/icons-material/CheckCircleOutline';

type SuccessProps = {
  message: string;
  title?: string;
  onDismiss?: () => void;
  autoDismiss?: boolean;
  autoDismissDelay?: number;
  fullWidth?: boolean;
};

export default function Success({
  message,
  title = 'Success',
  onDismiss,
  autoDismiss = false,
  autoDismissDelay = 5000,
  fullWidth = false,
}: SuccessProps) {
  useEffect(() => {
    if (autoDismiss && onDismiss) {
      const timer = setTimeout(() => {
        onDismiss();
      }, autoDismissDelay);

      return () => clearTimeout(timer);
    }
  }, [autoDismiss, autoDismissDelay, onDismiss]);

  return (
    <Box sx={{ width: fullWidth ? '100%' : 'auto', p: 2 }}>
      <Alert
        severity="success"
        icon={<CheckCircleOutlineIcon />}
        onClose={onDismiss}
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

