import React from 'react';
import { Box, IconButton, Button } from '@mui/material';
import ChevronLeftIcon from '@mui/icons-material/ChevronLeft';
import ChevronRightIcon from '@mui/icons-material/ChevronRight';

interface PaginationProps {
  currentPage: number;
  pageSize: number;
  totalItems: number;
  onPageChange?: (page: number) => void;
}

export default function Pagination({ currentPage, pageSize, totalItems, onPageChange }: PaginationProps) {
  const totalPages = Math.max(1, Math.ceil(totalItems / pageSize));
  const handlePrevious = () => {
    if (currentPage > 1 && onPageChange) {
      onPageChange(currentPage - 1);
    }
  };

  const handleNext = () => {
    if (currentPage < totalPages && onPageChange) {
      onPageChange(currentPage + 1);
    }
  };

  const handlePageClick = (page: number) => {
    if (onPageChange) {
      onPageChange(page);
    }
  };

  if (totalPages <= 1) {
    return null;
  }

  return (
    <Box
      sx={{
        display: 'flex',
        flexDirection: 'row',
        alignItems: 'center',
        justifyContent: 'center',
        gap: 1.5,
        px: 1,
      }}
    >
      <IconButton
        onClick={handlePrevious}
        disabled={currentPage === 1}
        sx={{
          width: 40,
          height: 40,
          border: '1px solid',
          borderColor: 'primary.main',
          color: 'text.primary',
          fontSize: '1.5rem',
          '&:hover': {
            borderWidth: 2,
            transform: 'scale(1.05)',
          },
          transition: 'all 0.075s',
        }}
      >
        <ChevronLeftIcon sx={{ fontSize: 20 }} />
      </IconButton>

      {Array.from({ length: totalPages }, (_, i) => i + 1).map((page) => (
        <Button
          key={page}
          onClick={() => handlePageClick(page)}
          variant={currentPage === page ? 'contained' : 'outlined'}
          sx={{
            minWidth: 36,
            height: 36,
            borderRadius: '50%',
            fontSize: '0.9rem',
            border: '1px solid',
            borderColor: 'primary.main',
            color: currentPage === page ? 'primary.contrastText' : 'text.primary',
            bgcolor: currentPage === page ? 'primary.main' : 'transparent',
            '&:hover': {
              borderWidth: 2,
              transform: 'scale(1.05)',
            },
            transition: 'all 0.075s',
          }}
        >
          {page}
        </Button>
      ))}

      <IconButton
        onClick={handleNext}
        disabled={currentPage === totalPages}
        sx={{
          width: 40,
          height: 40,
          border: '1px solid',
          borderColor: 'primary.main',
          color: 'text.primary',
          fontSize: '1.5rem',
          '&:hover': {
            borderWidth: 2,
            transform: 'scale(1.05)',
          },
          transition: 'all 0.075s',
        }}
      >
        <ChevronRightIcon sx={{ fontSize: 20 }} />
      </IconButton>
    </Box>
  );
}

