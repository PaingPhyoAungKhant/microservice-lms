import React, { useState, useCallback } from 'react';
import {
  Box,
  Typography,
  LinearProgress,
  Paper,
  IconButton,
} from '@mui/material';
import CloudUploadIcon from '@mui/icons-material/CloudUpload';
import CloseIcon from '@mui/icons-material/Close';
import ImageIcon from '@mui/icons-material/Image';
import { useUploadFileMutation } from '../../../infrastructure/api/rtk/filesApi';
import { getErrorMessage } from '../../../infrastructure/api/utils';
import Button from '../common/Button';

interface FileUploaderProps {
  bucketName: string;
  onUploadComplete: (fileId: string) => void;
  onClose?: () => void;
  acceptedTypes?: string[];
  maxSizeMB?: number;
}

export default function FileUploader({
  bucketName,
  onUploadComplete,
  onClose,
  acceptedTypes = ['image/*'],
  maxSizeMB = 10,
}: FileUploaderProps) {
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [preview, setPreview] = useState<string | null>(null);
  const [dragActive, setDragActive] = useState(false);
  const [uploadFile, { isLoading, error }] = useUploadFileMutation();

  const handleFileSelect = useCallback((file: File) => {
    if (acceptedTypes.length > 0) {
      const isAccepted = acceptedTypes.some((type) => {
        if (type.endsWith('/*')) {
          const baseType = type.split('/')[0];
          return file.type.startsWith(baseType + '/');
        }
        return file.type === type;
      });
      if (!isAccepted) {
        alert(`File type not accepted. Accepted types: ${acceptedTypes.join(', ')}`);
        return;
      }
    }

    const maxSizeBytes = maxSizeMB * 1024 * 1024;
    if (file.size > maxSizeBytes) {
      alert(`File size exceeds ${maxSizeMB}MB limit`);
      return;
    }

    setSelectedFile(file);

    if (file.type.startsWith('image/')) {
      const reader = new FileReader();
      reader.onloadend = () => {
        setPreview(reader.result as string);
      };
      reader.readAsDataURL(file);
    } else {
      setPreview(null);
    }
  }, [acceptedTypes, maxSizeMB]);

  const handleDrag = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    if (e.type === 'dragenter' || e.type === 'dragover') {
      setDragActive(true);
    } else if (e.type === 'dragleave') {
      setDragActive(false);
    }
  }, []);

  const handleDrop = useCallback(
    (e: React.DragEvent) => {
      e.preventDefault();
      e.stopPropagation();
      setDragActive(false);

      if (e.dataTransfer.files && e.dataTransfer.files[0]) {
        handleFileSelect(e.dataTransfer.files[0]);
      }
    },
    [handleFileSelect]
  );

  const handleFileInputChange = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      if (e.target.files && e.target.files[0]) {
        handleFileSelect(e.target.files[0]);
      }
    },
    [handleFileSelect]
  );

  const handleUpload = useCallback(async () => {
    if (!selectedFile) return;

    try {
      const result = await uploadFile({
        file: selectedFile,
        bucketName,
        tags: ['uploaded'],
      }).unwrap();

      onUploadComplete(result.id);
      setSelectedFile(null);
      setPreview(null);
    } catch (err) {
      const errorMessage = getErrorMessage(err, 'Upload failed');
      console.error('Upload failed:', {
        error: err,
        message: errorMessage,
        fileName: selectedFile.name,
        fileSize: selectedFile.size,
        fileType: selectedFile.type,
        bucketName,
      });
    }
  }, [selectedFile, bucketName, uploadFile, onUploadComplete]);

  const handleClear = useCallback(() => {
    setSelectedFile(null);
    setPreview(null);
  }, []);

  return (
    <Paper
      sx={{
        p: 3,
        border: dragActive ? '2px dashed' : '2px dashed',
        borderColor: dragActive ? 'primary.main' : 'grey.300',
        backgroundColor: dragActive ? 'action.hover' : 'background.paper',
        transition: 'all 0.2s',
      }}
      onDragEnter={handleDrag}
      onDragLeave={handleDrag}
      onDragOver={handleDrag}
      onDrop={handleDrop}
    >
      {onClose && (
        <Box sx={{ display: 'flex', justifyContent: 'flex-end', mb: 2 }}>
          <IconButton onClick={onClose} size="small">
            <CloseIcon />
          </IconButton>
        </Box>
      )}

      {!selectedFile ? (
        <Box
          sx={{
            display: 'flex',
            flexDirection: 'column',
            alignItems: 'center',
            justifyContent: 'center',
            py: 4,
            cursor: 'pointer',
          }}
          onClick={() => document.getElementById('file-input')?.click()}
        >
          <CloudUploadIcon sx={{ fontSize: 48, color: 'text.secondary', mb: 2 }} />
          <Typography variant="h6" gutterBottom>
            Drag & drop a file here
          </Typography>
          <Typography variant="body2" color="text.secondary" gutterBottom>
            or click to browse
          </Typography>
          <Typography variant="caption" color="text.secondary" sx={{ mt: 1 }}>
            Accepted types: {acceptedTypes.join(', ')} | Max size: {maxSizeMB}MB
          </Typography>
          <input
            id="file-input"
            type="file"
            accept={acceptedTypes.join(',')}
            onChange={handleFileInputChange}
            style={{ display: 'none' }}
          />
        </Box>
      ) : (
        <Box>
          {preview ? (
            <Box sx={{ mb: 2, textAlign: 'center' }}>
              <img
                src={preview}
                alt="Preview"
                style={{
                  maxWidth: '100%',
                  maxHeight: '300px',
                  objectFit: 'contain',
                  borderRadius: '8px',
                }}
              />
            </Box>
          ) : (
            <Box sx={{ display: 'flex', alignItems: 'center', mb: 2, p: 2, bgcolor: 'grey.100', borderRadius: 1 }}>
              <ImageIcon sx={{ mr: 2, color: 'text.secondary' }} />
              <Box sx={{ flex: 1 }}>
                <Typography variant="body1">{selectedFile.name}</Typography>
                <Typography variant="caption" color="text.secondary">
                  {(selectedFile.size / 1024 / 1024).toFixed(2)} MB
                </Typography>
              </Box>
            </Box>
          )}

          {isLoading && (
            <Box sx={{ mb: 2 }}>
              <LinearProgress />
              <Typography variant="caption" color="text.secondary" sx={{ mt: 1, display: 'block' }}>
                Uploading...
              </Typography>
            </Box>
          )}

          {error && (
            <Typography variant="body2" color="error" sx={{ mb: 2 }}>
              {getErrorMessage(error, 'Upload failed. Please try again.')}
            </Typography>
          )}

          <Box sx={{ display: 'flex', gap: 2, justifyContent: 'flex-end' }}>
            <Button variant="outline" onClick={handleClear} disabled={isLoading}>
              Clear
            </Button>
            <Button variant="fill" onClick={handleUpload} disabled={isLoading}>
              {isLoading ? 'Uploading...' : 'Upload'}
            </Button>
          </Box>
        </Box>
      )}
    </Paper>
  );
}

