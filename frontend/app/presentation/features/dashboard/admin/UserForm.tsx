import React, { useState, useEffect } from 'react';
import {
  Box,
  Typography,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Select,
  MenuItem,
  FormControl,
  InputLabel,
} from '@mui/material';
import InputField from '../../../components/forms/InputField';
import Button from '../../../components/common/Button';
import type { User, UserRole, UserStatus } from '../../../../domain/entities/User';

interface UserFormProps {
  open: boolean;
  user?: User | null;
  onClose: () => void;
  onSubmit: (data: {
    email?: string;
    username?: string;
    password?: string;
    role?: UserRole;
    status?: UserStatus;
  }) => Promise<void>;
  loading?: boolean;
}

export default function UserForm({ open, user, onClose, onSubmit, loading = false }: UserFormProps) {
  const isEdit = !!user;
  const [formData, setFormData] = useState({
    email: '',
    username: '',
    password: '',
    role: 'student' as UserRole,
    status: 'active' as UserStatus,
  });

  useEffect(() => {
    if (user) {
      setFormData({
        email: user.email || '',
        username: user.username || '',
        password: '',
        role: user.role,
        status: user.status || 'active',
      });
    } else {
      setFormData({
        email: '',
        username: '',
        password: '',
        role: 'student',
        status: 'active',
      });
    }
  }, [user, open]);

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
      await onSubmit(isEdit ? { ...formData, password: undefined } : formData);
      onClose();
    } catch (err) {
      console.error(err);
    }
  };

  return (
    <Dialog open={open} onClose={onClose} maxWidth="sm" fullWidth>
      <DialogTitle>{isEdit ? 'Edit User' : 'Create User'}</DialogTitle>
      <form onSubmit={handleSubmit}>
        <DialogContent>
          <InputField
            label="Username"
            name="username"
            value={formData.username}
            onChange={handleChange}
            required
            disabled={loading}
          />
          <InputField
            label="Email"
            type="email"
            name="email"
            value={formData.email}
            onChange={handleChange}
            required
            disabled={loading}
          />
          {!isEdit && (
            <InputField
              label="Password"
              type="password"
              name="password"
              value={formData.password}
              onChange={handleChange}
              required
              disabled={loading}
              helperText="Password must be at least 8 characters"
            />
          )}
          <Box sx={{ display: 'flex', gap: 2, mt: 2 }}>
            <FormControl fullWidth sx={{ mb: 2 }}>
              <InputLabel>Role</InputLabel>
              <Select
                name="role"
                value={formData.role}
                onChange={(e) =>
                  setFormData((prev) => ({ ...prev, role: e.target.value as UserRole }))
                }
                label="Role"
                disabled={loading}
              >
                <MenuItem value="student">Student</MenuItem>
                <MenuItem value="instructor">Instructor</MenuItem>
                <MenuItem value="admin">Admin</MenuItem>
              </Select>
            </FormControl>
            <FormControl fullWidth sx={{ mb: 2 }}>
              <InputLabel>Status</InputLabel>
              <Select
                name="status"
                value={formData.status}
                onChange={(e) =>
                  setFormData((prev) => ({ ...prev, status: e.target.value as UserStatus }))
                }
                label="Status"
                disabled={loading}
              >
                <MenuItem value="active">Active</MenuItem>
                <MenuItem value="inactive">Inactive</MenuItem>
                <MenuItem value="pending">Pending</MenuItem>
                <MenuItem value="banned">Banned</MenuItem>
              </Select>
            </FormControl>
          </Box>
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

