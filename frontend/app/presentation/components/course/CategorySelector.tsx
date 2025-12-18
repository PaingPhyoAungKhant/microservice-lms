import React, { useState, useMemo } from 'react';
import {
  Box,
  Typography,
  Chip,
  Autocomplete,
  TextField,
  Paper,
} from '@mui/material';
import AddIcon from '@mui/icons-material/Add';
import CloseIcon from '@mui/icons-material/Close';
import { useGetCategoriesQuery } from '../../../infrastructure/api/rtk/categoriesApi';
import Button from '../common/Button';
import type { Category } from '../../../domain/entities/Category';

interface CategorySelectorProps {
  selectedCategories: Category[];
  onChange: (categories: Category[]) => void;
  error?: string;
}

export default function CategorySelector({
  selectedCategories,
  onChange,
  error,
}: CategorySelectorProps) {
  const [searchQuery, setSearchQuery] = useState('');
  const { data, isLoading } = useGetCategoriesQuery({
    search: searchQuery || undefined,
    limit: 50,
  });

  const availableCategories = useMemo(() => {
    if (!data) return [];
    const selectedIds = new Set(selectedCategories.map((cat) => cat.id));
    return data.categories.filter((cat) => !selectedIds.has(cat.id));
  }, [data, selectedCategories]);

  const handleAddCategory = (category: Category | null) => {
    if (!category) return;
    if (!selectedCategories.find((cat) => cat.id === category.id)) {
      onChange([...selectedCategories, category]);
    }
    setSearchQuery('');
  };

  const handleRemoveCategory = (categoryId: string) => {
    onChange(selectedCategories.filter((cat) => cat.id !== categoryId));
  };

  return (
    <Box>
      <Typography variant="subtitle2" gutterBottom sx={{ fontWeight: 600 }}>
        Categories
      </Typography>

      {selectedCategories.length > 0 && (
        <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 1, mb: 2 }}>
          {selectedCategories.map((category) => (
            <Chip
              key={category.id}
              label={category.name}
              onDelete={() => handleRemoveCategory(category.id)}
              deleteIcon={<CloseIcon />}
              color="primary"
              variant="filled"
            />
          ))}
        </Box>
      )}

      <Autocomplete
        options={availableCategories}
        getOptionLabel={(option) => option.name}
        loading={isLoading}
        inputValue={searchQuery}
        onInputChange={(_e, newValue) => setSearchQuery(newValue)}
        onChange={(_e, value) => handleAddCategory(value)}
        renderInput={(params) => (
          <TextField
            {...params}
            placeholder="Search and add categories..."
            variant="outlined"
            error={!!error}
            helperText={error}
            InputProps={{
              ...params.InputProps,
              endAdornment: (
                <>
                  {isLoading ? <Box sx={{ width: 20, height: 20 }} /> : null}
                  {params.InputProps.endAdornment}
                </>
              ),
            }}
          />
        )}
        renderOption={(props, option) => (
          <Box component="li" {...props} key={option.id}>
            <Box>
              <Typography variant="body1">{option.name}</Typography>
              {option.description && (
                <Typography variant="caption" color="text.secondary">
                  {option.description}
                </Typography>
              )}
            </Box>
          </Box>
        )}
        PaperComponent={({ children, ...other }) => (
          <Paper {...other}>{children}</Paper>
        )}
        noOptionsText={
          searchQuery
            ? `No categories found matching "${searchQuery}"`
            : 'No categories available'
        }
        sx={{ mb: 1 }}
      />

      {selectedCategories.length === 0 && (
        <Typography variant="caption" color="text.secondary" sx={{ mt: 1, display: 'block' }}>
          No categories selected. Add categories to organize this course.
        </Typography>
      )}
    </Box>
  );
}

