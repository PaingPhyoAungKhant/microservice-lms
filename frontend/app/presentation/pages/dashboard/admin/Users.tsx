import React, { useState, useEffect, useCallback } from 'react';
import {
  Box,
  Typography,
  IconButton,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  DialogContentText,
  TablePagination,
  Card,
  CardContent,
} from '@mui/material';
import AddIcon from '@mui/icons-material/Add';
import { useUsers, useCreateUser, useUpdateUser, useDeleteUser } from '../../../hooks/useUsers';
import { useDashboardAuth } from '../../../hooks/useDashboardAuth';
import DashboardLayout from '../../../components/dashboard/Layout';
import UserList from '../../../features/dashboard/admin/UserList';
import UserFilters from '../../../features/dashboard/admin/UserFilters';
import UserForm from '../../../features/dashboard/admin/UserForm';
import Loading from '../../../components/common/Loading';
import Error from '../../../components/common/Error';
import Success from '../../../components/common/Success';
import Button from '../../../components/common/Button';
import type { User, UserRole, UserStatus } from '../../../../domain/entities/User';

const ITEMS_PER_PAGE = 10;

export default function AdminUsers() {
  const { user: currentUser, logout } = useDashboardAuth();
  
  const [page, setPage] = useState(0);
  const [searchQuery, setSearchQuery] = useState('');
  const [roleFilter, setRoleFilter] = useState<UserRole | ''>('');
  const [statusFilter, setStatusFilter] = useState<UserStatus | ''>('');
  const [formOpen, setFormOpen] = useState(false);
  const [editingUser, setEditingUser] = useState<User | null>(null);
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [userToDelete, setUserToDelete] = useState<User | null>(null);
  const [successMessage, setSuccessMessage] = useState<string | null>(null);

  const { users, total, loading, error, refetch } = useUsers({
    searchQuery: searchQuery || undefined,
    role: roleFilter || undefined,
    status: statusFilter || undefined,
    limit: ITEMS_PER_PAGE,
    offset: page * ITEMS_PER_PAGE,
    sortColumn: 'created_at',
    sortDirection: 'desc',
  });

  const { loading: creating, error: createError, createUser, reset: resetCreate } = useCreateUser();
  const { loading: updating, error: updateError, updateUser, reset: resetUpdate } = useUpdateUser();
  const { loading: deleting, error: deleteError, deleteUser, reset: resetDelete } = useDeleteUser();

  const handleCreate = () => {
    setEditingUser(null);
    setFormOpen(true);
  };

  const handleEdit = (user: User) => {
    setEditingUser(user);
    setFormOpen(true);
  };

  const handleDelete = (user: User) => {
    setUserToDelete(user);
    setDeleteDialogOpen(true);
  };

  const handleFormSubmit = async (data: {
    email?: string;
    username?: string;
    password?: string;
    role?: UserRole;
    status?: UserStatus;
  }) => {
    try {
      if (editingUser) {
        await updateUser(editingUser.id, data);
        setSuccessMessage('User updated successfully');
      } else {
        if (!data.email || !data.username || !data.password || !data.role) {
          throw new Error('Missing required fields');
        }
        await createUser({
          email: data.email,
          username: data.username,
          password: data.password,
          role: data.role,
        });
        setSuccessMessage('User created successfully');
      }
      setFormOpen(false);
      setEditingUser(null);
      refetch();
    } catch (err) {
      console.error(err);
    }
  };

  const handleConfirmDelete = async () => {
    if (!userToDelete) return;

    try {
      await deleteUser(userToDelete.id);
      setSuccessMessage('User deleted successfully');
      setDeleteDialogOpen(false);
      setUserToDelete(null);
      refetch();
    } catch (err) {
      console.error(err);
    }
  };

  const handlePageChange = (_event: unknown, newPage: number) => {
    setPage(newPage);
  };

  if (!currentUser || currentUser.role !== 'admin') {
    return (
      <Box sx={{ p: 4, textAlign: 'center' }}>
        <Typography variant="h5" color="error">
          Access Denied. Admin role required.
        </Typography>
      </Box>
    );
  }

  return (
    <DashboardLayout user={currentUser} onLogout={logout}>
      <Box>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
          <Typography variant="h4" component="h1" sx={{ fontWeight: 600 }}>
            User Management
          </Typography>
          <IconButton
            onClick={handleCreate}
            sx={{
              bgcolor: 'primary.main',
              color: 'primary.contrastText',
              '&:hover': { bgcolor: 'primary.dark' },
            }}
            aria-label="create user"
          >
            <AddIcon />
          </IconButton>
        </Box>

        {successMessage && (
          <Success
            message={successMessage}
            autoDismiss
            autoDismissDelay={3000}
            onDismiss={() => setSuccessMessage(null)}
            fullWidth
          />
        )}

        {(createError || updateError || deleteError) && (
          <Error
            message={
              createError?.message || updateError?.message || deleteError?.message || 'An error occurred'
            }
            onRetry={() => {
              resetCreate();
              resetUpdate();
              resetDelete();
            }}
            fullWidth
          />
        )}

        <Card>
          <CardContent>
            <UserFilters
              searchQuery={searchQuery}
              role={roleFilter}
              status={statusFilter}
              onSearchChange={setSearchQuery}
              onRoleChange={setRoleFilter}
              onStatusChange={setStatusFilter}
            />

            {error && (
              <Error
                message={error.message || 'Failed to load users'}
                onRetry={() => refetch()}
                fullWidth
              />
            )}

            {loading && <Loading variant="skeleton" fullWidth />}

            {!loading && !error && (
              <>
                <UserList
                  users={users}
                  loading={loading}
                  onEdit={handleEdit}
                  onDelete={handleDelete}
                  page={page}
                  itemsPerPage={ITEMS_PER_PAGE}
                />
                <TablePagination
                  component="div"
                  count={total}
                  page={page}
                  onPageChange={handlePageChange}
                  rowsPerPage={ITEMS_PER_PAGE}
                  rowsPerPageOptions={[]}
                />
              </>
            )}
          </CardContent>
        </Card>

        <UserForm
          open={formOpen}
          user={editingUser}
          onClose={() => {
            setFormOpen(false);
            setEditingUser(null);
            resetCreate();
            resetUpdate();
          }}
          onSubmit={handleFormSubmit}
          loading={creating || updating}
        />

        <Dialog open={deleteDialogOpen} onClose={() => setDeleteDialogOpen(false)}>
          <DialogTitle>Delete User</DialogTitle>
          <DialogContent>
            <DialogContentText>
              Are you sure you want to delete user &quot;{userToDelete?.username}&quot;? This action
              cannot be undone.
            </DialogContentText>
          </DialogContent>
          <DialogActions>
            <Button onClick={() => setDeleteDialogOpen(false)} variant="outline" disabled={deleting}>
              Cancel
            </Button>
            <Button
              onClick={handleConfirmDelete}
              variant="fill"
              color="error"
              disabled={deleting}
            >
              {deleting ? 'Deleting...' : 'Delete'}
            </Button>
          </DialogActions>
        </Dialog>
      </Box>
    </DashboardLayout>
  );
}

