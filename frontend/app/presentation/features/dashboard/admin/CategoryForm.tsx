import React, { useState, useEffect } from 'react';
import {
  Box,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
} from '@mui/material';
import InputField from '../../../components/forms/InputField';
import Button from '../../../components/common/Button';
import type { Category } from '../../../../domain/entities/Category';

interface CategoryFormProps {
  open: boolean;
  category?: Category | null;
  onClose: () => void;
  onSubmit: (data: {
    name: string;
    description?: string;
  }) => Promise<void>;
  loading?: boolean;
}

export default function CategoryForm({ open, category, onClose, onSubmit, loading = false }: CategoryFormProps) {
  const isEdit = !!category;
  const [formData, setFormData] = useState({
    name: '',
    description: '',
  });

  useEffect(() => {
    if (category) {
      setFormData({
        name: category.name || '',
        description: category.description || '',
      });
    } else {
      setFormData({
        name: '',
        description: '',
      });
    }
  }, [category, open]);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData((prev) => ({
      ...prev,
      [name]: value,
    }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await onSubmit({
        name: formData.name,
        description: formData.description || undefined,
      });
      onClose();
    } catch (err) {
      console.error(err);
    }
  };

  return (
    <Dialog open={open} onClose={onClose} maxWidth="sm" fullWidth>
      <DialogTitle>{isEdit ? 'Edit Category' : 'Create Category'}</DialogTitle>
      <form onSubmit={handleSubmit}>
        <DialogContent>
          <InputField
            label="Name"
            name="name"
            value={formData.name}
            onChange={handleChange}
            required
            disabled={loading}
            sx={{ mb: 2 }}
          />
          <InputField
            label="Description"
            name="description"
            value={formData.description}
            onChange={handleChange}
            disabled={loading}
            multiline
            rows={4}
            sx={{ mb: 2 }}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={onClose} variant="outline" disabled={loading}>
            Cancel
          </Button>
          <Button type="submit" variant="fill" color="primary" disabled={loading}>
            {loading ? (isEdit ? 'Updating...' : 'Creating...') : isEdit ? 'Update' : 'Create'}
          </Button>
        </DialogActions>
      </form>
    </Dialog>
  );
}

