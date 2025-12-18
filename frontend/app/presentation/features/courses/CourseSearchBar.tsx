import React, { useState } from 'react';
import { Box, TextField, IconButton, InputAdornment } from '@mui/material';
import SearchIcon from '@mui/icons-material/Search';

interface CourseSearchBarProps {
  onSearch?: (query: string) => void;
  placeholder?: string;
}

export default function CourseSearchBar({ onSearch, placeholder = 'Search ...' }: CourseSearchBarProps) {
  const [searchQuery, setSearchQuery] = useState('');

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (onSearch) {
      onSearch(searchQuery);
    }
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setSearchQuery(e.target.value);
  };

  return (
    <Box
      component="form"
      onSubmit={handleSubmit}
      sx={{
        display: 'flex',
        alignItems: 'center',
        width: '100%',
        height: 64,
        bgcolor: 'grey.200',
        borderRadius: '9999px',
        px: 3,
      }}
    >
      <TextField
        fullWidth
        placeholder={placeholder}
        value={searchQuery}
        onChange={handleChange}
        variant="standard"
        InputProps={{
          disableUnderline: true,
          endAdornment: (
            <InputAdornment position="end">
              <IconButton
                type="submit"
                sx={{
                  transform: 'scale(1.3)',
                  '&:hover': { transform: 'scale(1.45)' },
                  transition: 'all 0.3s ease-in-out',
                }}
              >
                <SearchIcon sx={{ fontSize: 32 }} />
              </IconButton>
            </InputAdornment>
          ),
        }}
        sx={{
          '& .MuiInputBase-input': {
            fontSize: '1.5rem',
            py: 2,
            '&::placeholder': {
              color: 'text.primary',
              opacity: 1,
            },
          },
        }}
      />
    </Box>
  );
}

