import React from 'react';
import { Box, TextField, MenuItem, Grid } from '@mui/material';
import type { UserRole, UserStatus } from '../../../../domain/entities/User';

interface UserFiltersProps {
  searchQuery: string;
  role: UserRole | '';
  status: UserStatus | '';
  onSearchChange: (value: string) => void;
  onRoleChange: (value: UserRole | '') => void;
  onStatusChange: (value: UserStatus | '') => void;
}

export default function UserFilters({
  searchQuery,
  role,
  status,
  onSearchChange,
  onRoleChange,
  onStatusChange,
}: UserFiltersProps) {
  return (
    <Box sx={{ mb: 3 }}>
      <Grid container spacing={2}>
        <Grid item xs={12} md={4}>
          <TextField
            fullWidth
            label="Search"
            placeholder="Search by username or email"
            value={searchQuery}
            onChange={(e) => onSearchChange(e.target.value)}
            variant="outlined"
          />
        </Grid>
        <Grid item xs={12} sm={6} md={4}>
          <TextField
            fullWidth
            select
            label="Role"
            value={role}
            onChange={(e) => onRoleChange(e.target.value as UserRole | '')}
            variant="outlined"
            sx={{ minWidth: 180 }}
          >
            <MenuItem value="">All Roles</MenuItem>
            <MenuItem value="student">Student</MenuItem>
            <MenuItem value="instructor">Instructor</MenuItem>
            <MenuItem value="admin">Admin</MenuItem>
          </TextField>
        </Grid>
        <Grid item xs={12} sm={6} md={4}>
          <TextField
            fullWidth
            select
            label="Status"
            value={status}
            onChange={(e) => onStatusChange(e.target.value as UserStatus | '')}
            variant="outlined"
            sx={{ minWidth: 180 }}
          >
            <MenuItem value="">All Statuses</MenuItem>
            <MenuItem value="active">Active</MenuItem>
            <MenuItem value="inactive">Inactive</MenuItem>
            <MenuItem value="pending">Pending</MenuItem>
            <MenuItem value="banned">Banned</MenuItem>
          </TextField>
        </Grid>
      </Grid>
    </Box>
  );
}

