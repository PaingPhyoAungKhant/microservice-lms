import React from 'react';
import { Box, TextField, Grid } from '@mui/material';

interface CategoryFiltersProps {
  searchQuery: string;
  onSearchChange: (value: string) => void;
}

export default function CategoryFilters({
  searchQuery,
  onSearchChange,
}: CategoryFiltersProps) {
  return (
    <Box sx={{ mb: 3 }}>
      <Grid container spacing={2}>
        <Grid item xs={12} md={4}>
          <TextField
            fullWidth
            label="Search"
            placeholder="Search by name or description"
            value={searchQuery}
            onChange={(e) => onSearchChange(e.target.value)}
            variant="outlined"
          />
        </Grid>
      </Grid>
    </Box>
  );
}

