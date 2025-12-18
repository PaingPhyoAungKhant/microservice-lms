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
import { useCategories, useCreateCategory, useUpdateCategory, useDeleteCategory } from '../../../hooks/useCategories';
import { useDashboardAuth } from '../../../hooks/useDashboardAuth';
import DashboardLayout from '../../../components/dashboard/Layout';
import CategoryList from '../../../features/dashboard/admin/CategoryList';
import CategoryFilters from '../../../features/dashboard/admin/CategoryFilters';
import CategoryForm from '../../../features/dashboard/admin/CategoryForm';
import Loading from '../../../components/common/Loading';
import Error from '../../../components/common/Error';
import Success from '../../../components/common/Success';
import Button from '../../../components/common/Button';
import type { Category } from '../../../../domain/entities/Category';

const ITEMS_PER_PAGE = 10;

export default function AdminCategories() {
  const { user: currentUser, logout } = useDashboardAuth();
  
  const [page, setPage] = useState(0);
  const [searchQuery, setSearchQuery] = useState('');
  const [formOpen, setFormOpen] = useState(false);
  const [editingCategory, setEditingCategory] = useState<Category | null>(null);
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [categoryToDelete, setCategoryToDelete] = useState<Category | null>(null);
  const [successMessage, setSuccessMessage] = useState<string | null>(null);

  const { categories, total, loading, error, refetch } = useCategories({
    search: searchQuery || undefined,
    limit: ITEMS_PER_PAGE,
    offset: page * ITEMS_PER_PAGE,
    sortColumn: 'created_at',
    sortDirection: 'desc',
  });

  const { loading: creating, error: createError, createCategory, reset: resetCreate } = useCreateCategory();
  const { loading: updating, error: updateError, updateCategory, reset: resetUpdate } = useUpdateCategory();
  const { loading: deleting, error: deleteError, deleteCategory, reset: resetDelete } = useDeleteCategory();

  const handleCreate = () => {
    setEditingCategory(null);
    setFormOpen(true);
  };

  const handleEdit = (category: Category) => {
    setEditingCategory(category);
    setFormOpen(true);
  };

  const handleDelete = (category: Category) => {
    setCategoryToDelete(category);
    setDeleteDialogOpen(true);
  };

  const handleFormSubmit = async (data: {
    name: string;
    description?: string;
  }) => {
    try {
      if (editingCategory) {
        await updateCategory(editingCategory.id, data);
        setSuccessMessage('Category updated successfully');
      } else {
        await createCategory({
          name: data.name,
          description: data.description,
        });
        setSuccessMessage('Category created successfully');
      }
      setFormOpen(false);
      setEditingCategory(null);
      refetch();
    } catch (err) {
      console.error(err);
    }
  };

  const handleConfirmDelete = async () => {
    if (!categoryToDelete) return;

    try {
      await deleteCategory(categoryToDelete.id);
      setSuccessMessage('Category deleted successfully');
      setDeleteDialogOpen(false);
      setCategoryToDelete(null);
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
            Category Management
          </Typography>
          <IconButton
            onClick={handleCreate}
            sx={{
              bgcolor: 'primary.main',
              color: 'primary.contrastText',
              '&:hover': { bgcolor: 'primary.dark' },
            }}
            aria-label="create category"
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
            <CategoryFilters
              searchQuery={searchQuery}
              onSearchChange={setSearchQuery}
            />

            {error && (
              <Error
                message={error.message || 'Failed to load categories'}
                onRetry={() => refetch()}
                fullWidth
              />
            )}

            {loading && <Loading variant="skeleton" fullWidth />}

            {!loading && !error && (
              <>
                <CategoryList
                  categories={categories}
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

        <CategoryForm
          open={formOpen}
          category={editingCategory}
          onClose={() => {
            setFormOpen(false);
            setEditingCategory(null);
            resetCreate();
            resetUpdate();
          }}
          onSubmit={handleFormSubmit}
          loading={creating || updating}
        />

        <Dialog open={deleteDialogOpen} onClose={() => setDeleteDialogOpen(false)}>
          <DialogTitle>Delete Category</DialogTitle>
          <DialogContent>
            <DialogContentText>
              Are you sure you want to delete category &quot;{categoryToDelete?.name}&quot;? This action
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

