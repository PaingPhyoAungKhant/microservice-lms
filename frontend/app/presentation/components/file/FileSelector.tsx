import React, { useState, useCallback, useEffect } from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Box,
  Typography,
  Card,
  CardMedia,
  CardContent,
  IconButton,
  TextField,
  InputAdornment,
  CircularProgress,
  Pagination,
  Paper,
} from '@mui/material';
import SearchIcon from '@mui/icons-material/Search';
import CloseIcon from '@mui/icons-material/Close';
import CheckCircleIcon from '@mui/icons-material/CheckCircle';
import CloudUploadIcon from '@mui/icons-material/CloudUpload';
import ImageIcon from '@mui/icons-material/Image';
import BrokenImageIcon from '@mui/icons-material/BrokenImage';
import { useListFilesQuery } from '../../../infrastructure/api/rtk/filesApi';
import { getFileDownloadUrl } from '../../../infrastructure/api/utils';
import Button from '../common/Button';
import FileUploader from './FileUploader';
import type { File } from '../../../domain/entities/File';

interface FileSelectorProps {
  open: boolean;
  bucketName: string;
  onClose: () => void;
  onSelect: (fileId: string) => void;
  selectedFileId?: string | null;
  acceptedTypes?: string[];
}

const ITEMS_PER_PAGE = 12;

export default function FileSelector({
  open,
  bucketName,
  onClose,
  onSelect,
  selectedFileId,
  acceptedTypes = ['image/*'],
}: FileSelectorProps) {
  const [page, setPage] = useState(1);
  const [searchQuery, setSearchQuery] = useState('');
  const [showUploader, setShowUploader] = useState(false);
  const [uploadedFileId, setUploadedFileId] = useState<string | null>(null);
  const [imageErrors, setImageErrors] = useState<Set<string>>(new Set());

  const { data, isLoading, error, refetch } = useListFilesQuery({
    bucketName,
    limit: ITEMS_PER_PAGE,
    offset: (page - 1) * ITEMS_PER_PAGE,
  });

  useEffect(() => {
    if (uploadedFileId) {
      refetch();
      setUploadedFileId(null);
    }
  }, [uploadedFileId, refetch]);

  const handleSelect = useCallback(
    (fileId: string) => {
      onSelect(fileId);
      onClose();
    },
    [onSelect, onClose]
  );

  const handleUploadComplete = useCallback((fileId: string) => {
    setUploadedFileId(fileId);
    setShowUploader(false);
    handleSelect(fileId);
  }, [handleSelect]);

  const handleImageError = useCallback((fileId: string) => {
    setImageErrors((prev) => new Set(prev).add(fileId));
  }, []);

  const getFileUrl = useCallback((file: File) => {
    return getFileDownloadUrl(bucketName, file.id);
  }, [bucketName]);

  const totalPages = data ? Math.ceil(data.total / ITEMS_PER_PAGE) : 1;

  return (
    <Dialog open={open} onClose={onClose} maxWidth="lg" fullWidth>
      <DialogTitle>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <Typography variant="h6">Select File from {bucketName}</Typography>
          <IconButton onClick={onClose} size="small">
            <CloseIcon />
          </IconButton>
        </Box>
      </DialogTitle>
      <DialogContent>
        {showUploader ? (
          <FileUploader
            bucketName={bucketName}
            onUploadComplete={handleUploadComplete}
            onClose={() => setShowUploader(false)}
            acceptedTypes={acceptedTypes}
          />
        ) : (
          <>
            <Box sx={{ display: 'flex', gap: 2, mb: 3 }}>
              <TextField
                fullWidth
                placeholder="Search files..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                InputProps={{
                  startAdornment: (
                    <InputAdornment position="start">
                      <SearchIcon />
                    </InputAdornment>
                  ),
                }}
              />
              <Button
                variant="fill"
                startIcon={<CloudUploadIcon />}
                onClick={() => setShowUploader(true)}
              >
                Upload New
              </Button>
            </Box>

            {error && (
              <Paper sx={{ p: 2, bgcolor: 'error.light', color: 'error.contrastText', mb: 2 }}>
                <Typography>Failed to load files. Please try again.</Typography>
              </Paper>
            )}

            {isLoading ? (
              <Box sx={{ display: 'flex', justifyContent: 'center', py: 4 }}>
                <CircularProgress />
              </Box>
            ) : data && data.files.length > 0 ? (
              <>
                <Box
                  sx={{
                    display: 'grid',
                    gridTemplateColumns: {
                      xs: '1fr',
                      sm: 'repeat(2, 1fr)',
                      md: 'repeat(3, 1fr)',
                    },
                    gap: 2,
                  }}
                >
                  {data.files
                    .filter((file) =>
                      searchQuery
                        ? file.originalFilename.toLowerCase().includes(searchQuery.toLowerCase())
                        : true
                    )
                    .map((file) => (
                      <Card
                        key={file.id}
                        sx={{
                          cursor: 'pointer',
                          position: 'relative',
                          border: selectedFileId === file.id ? '2px solid' : '1px solid',
                          borderColor: selectedFileId === file.id ? 'primary.main' : 'divider',
                          '&:hover': {
                            boxShadow: 3,
                          },
                        }}
                        onClick={() => handleSelect(file.id)}
                      >
                        {file.mimeType.startsWith('image/') && !imageErrors.has(file.id) ? (
                          <CardMedia
                            component="img"
                            height="200"
                            image={getFileUrl(file)}
                            alt={file.originalFilename}
                            sx={{ objectFit: 'cover' }}
                            onError={() => handleImageError(file.id)}
                          />
                        ) : (
                          <Box
                            sx={{
                              height: 200,
                              display: 'flex',
                              flexDirection: 'column',
                              alignItems: 'center',
                              justifyContent: 'center',
                              bgcolor: 'grey.100',
                              gap: 1,
                            }}
                          >
                            {file.mimeType.startsWith('image/') && imageErrors.has(file.id) ? (
                              <>
                                <BrokenImageIcon sx={{ fontSize: 48, color: 'text.secondary' }} />
                                <Typography variant="caption" color="text.secondary">
                                  Failed to load image
                                </Typography>
                              </>
                            ) : (
                              <>
                                <ImageIcon sx={{ fontSize: 48, color: 'text.secondary' }} />
                                <Typography variant="body2" color="text.secondary" sx={{ px: 2, textAlign: 'center' }}>
                              {file.originalFilename}
                            </Typography>
                              </>
                            )}
                          </Box>
                        )}
                        <CardContent>
                          <Typography variant="body2" noWrap title={file.originalFilename}>
                            {file.originalFilename}
                          </Typography>
                          <Typography variant="caption" color="text.secondary">
                            {(file.sizeBytes / 1024).toFixed(2)} KB
                          </Typography>
                        </CardContent>
                        {selectedFileId === file.id && (
                          <Box
                            sx={{
                              position: 'absolute',
                              top: 8,
                              right: 8,
                              bgcolor: 'primary.main',
                              borderRadius: '50%',
                              p: 0.5,
                            }}
                          >
                            <CheckCircleIcon sx={{ color: 'white', fontSize: 20 }} />
                          </Box>
                        )}
                      </Card>
                    ))}
                </Box>
                {totalPages > 1 && (
                  <Box sx={{ display: 'flex', justifyContent: 'center', mt: 3 }}>
                    <Pagination
                      count={totalPages}
                      page={page}
                      onChange={(_e, newPage) => setPage(newPage)}
                      color="primary"
                    />
                  </Box>
                )}
              </>
            ) : (
              <Box sx={{ textAlign: 'center', py: 4 }}>
                <Typography variant="body1" color="text.secondary">
                  No files found in this bucket.
                </Typography>
                <Button
                  variant="outline"
                  startIcon={<CloudUploadIcon />}
                  onClick={() => setShowUploader(true)}
                  sx={{ mt: 2 }}
                >
                  Upload First File
                </Button>
              </Box>
            )}
          </>
        )}
      </DialogContent>
      <DialogActions>
        <Button variant="outline" onClick={onClose}>
          Cancel
        </Button>
      </DialogActions>
    </Dialog>
  );
}

