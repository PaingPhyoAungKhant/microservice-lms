import React from 'react';
import { CircularProgress, Box, Skeleton, Stack } from '@mui/material';

type LoadingProps = {
  variant?: 'circular' | 'skeleton' | 'fullscreen';
  size?: number;
  message?: string;
  fullWidth?: boolean;
};

export default function Loading({
  variant = 'circular',
  size = 40,
  message,
  fullWidth = false,
}: LoadingProps) {
  if (variant === 'fullscreen') {
    return (
      <Box
        sx={{
          display: 'flex',
          flexDirection: 'column',
          alignItems: 'center',
          justifyContent: 'center',
          minHeight: '100vh',
          gap: 2,
        }}
      >
        <CircularProgress size={size} />
        {message && (
          <Box component="p" sx={{ color: 'text.secondary', mt: 2 }}>
            {message}
          </Box>
        )}
      </Box>
    );
  }

  if (variant === 'skeleton') {
    return (
      <Stack spacing={2} sx={{ width: fullWidth ? '100%' : 'auto' }}>
        <Skeleton variant="rectangular" width="100%" height={200} />
        <Skeleton variant="text" width="100%" height={40} />
        <Skeleton variant="text" width="80%" height={40} />
        <Skeleton variant="text" width="60%" height={40} />
      </Stack>
    );
  }

  return (
    <Box
      sx={{
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        justifyContent: 'center',
        p: 4,
        gap: 2,
        width: fullWidth ? '100%' : 'auto',
      }}
    >
      <CircularProgress size={size} />
      {message && (
        <Box component="p" sx={{ color: 'text.secondary' }}>
          {message}
        </Box>
      )}
    </Box>
  );
}

