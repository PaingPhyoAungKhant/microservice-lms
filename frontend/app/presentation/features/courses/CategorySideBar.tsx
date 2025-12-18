import React from 'react';
import { Box, Typography, Chip, CircularProgress } from '@mui/material';
import { useCategories } from '../../hooks/useCategories';

interface CategorySideBarProps {
  onCategorySelect?: (categoryId: string | undefined) => void;
  selectedCategory?: string;
}

export default function CategorySideBar({ onCategorySelect, selectedCategory }: CategorySideBarProps) {
  const { categories, loading, error } = useCategories({
    limit: 50,
    sortColumn: 'name',
    sortDirection: 'asc',
  });

  const handleAllClick = () => {
    onCategorySelect?.(undefined);
  };

  return (
    <Box sx={{ width: '100%' }}>
      <Typography variant="h4" sx={{ mb: 3, fontSize: '2.25rem', fontWeight: 600, color: 'text.secondary' }}>
        Filter Courses
      </Typography>

      <Box sx={{ mb: 3 }}>
        <Typography variant="h6" sx={{ mb: 2, fontSize: '1.5rem', fontWeight: 500, color: 'text.secondary' }}>
          Course Categories
        </Typography>
        {loading && (
          <Box sx={{ mt: 3, display: 'flex', justifyContent: 'center' }}>
            <CircularProgress size={24} />
          </Box>
        )}
        {error && (
          <Typography variant="body2" color="error" sx={{ mt: 2 }}>
            Failed to load categories
          </Typography>
        )}
        {!loading && !error && (
        <Box
          sx={{
            mt: 3,
              display: 'flex',
              flexDirection: { xs: 'row', md: 'column' },
              flexWrap: { xs: 'wrap', md: 'nowrap' },
              alignItems: { xs: 'flex-start', md: 'stretch' },
              gap: 1,
          }}
        >
            <Chip
              label="All"
              clickable
              onClick={handleAllClick}
              color={!selectedCategory ? 'primary' : 'default'}
              variant={!selectedCategory ? 'filled' : 'outlined'}
              sx={{ mr: { xs: 0, md: 0 }, width: { xs: 'auto', md: '100%' } }}
            />
          {categories.map((category) => (
              <Chip
                key={category.id}
                label={category.name}
                clickable
                onClick={() => onCategorySelect?.(category.id)}
                color={selectedCategory === category.id ? 'primary' : 'default'}
                variant={selectedCategory === category.id ? 'filled' : 'outlined'}
                sx={{ width: { xs: 'auto', md: '100%' } }}
              />
          ))}
        </Box>
        )}
      </Box>
    </Box>
  );
}

