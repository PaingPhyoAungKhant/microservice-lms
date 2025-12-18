import React, { useState, useRef } from 'react';
import {
  Box,
  Typography,
  Button,
  List,
  ListItem,
  ListItemText,
  IconButton,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  LinearProgress,
  Chip,
  Alert,
  Snackbar,
  CircularProgress,
} from '@mui/material';
import DeleteIcon from '@mui/icons-material/Delete';
import DownloadIcon from '@mui/icons-material/Download';
import UploadFileIcon from '@mui/icons-material/UploadFile';
import {
  useListZoomRecordingsQuery,
  useCreateZoomRecordingMutation,
  useDeleteZoomRecordingMutation,
} from '../../../infrastructure/api/rtk/zoomRecordingsApi';
import { useUploadFileMutation } from '../../../infrastructure/api/rtk/filesApi';
import { downloadFileWithAuth, getErrorMessage } from '../../../infrastructure/api/utils';

interface ZoomRecordingSectionProps {
  meetingId: string;
}

export default function ZoomRecordingSection({ meetingId }: ZoomRecordingSectionProps) {
  const { data: recordings = [], isLoading } = useListZoomRecordingsQuery(meetingId);
  const [createRecording, { isLoading: creating }] = useCreateZoomRecordingMutation();
  const [deleteRecording] = useDeleteZoomRecordingMutation();
  const [uploadFile, { isLoading: uploading }] = useUploadFileMutation();
  const [deleteDialogOpen, setDeleteDialogOpen] = useState<string | null>(null);
  const [uploadProgress, setUploadProgress] = useState(0);
  const [errorMessage, setErrorMessage] = useState<string | null>(null);
  const [downloadingFileId, setDownloadingFileId] = useState<string | null>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const handleUploadClick = () => {
    fileInputRef.current?.click();
  };

  const handleFileChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    setUploadProgress(0);
    try {
      const uploadedFile = await uploadFile({
        file,
        bucketName: 'zoom-recordings',
      }).unwrap();

      if (!meetingId || typeof meetingId !== 'string') {
        throw new Error('Meeting ID is required');
      }
      if (!uploadedFile?.id || typeof uploadedFile.id !== 'string') {
        throw new Error('File ID is required');
      }

      await createRecording({
        zoomMeetingId: meetingId,
        fileId: uploadedFile.id,
        fileSize: uploadedFile.sizeBytes,
      }).unwrap();

      if (fileInputRef.current) {
        fileInputRef.current.value = '';
      }
      setUploadProgress(0);
    } catch (err) {
      console.error('Failed to upload recording:', err);
      setUploadProgress(0);
      
      let errorMsg = getErrorMessage(err, 'Failed to upload recording');
      
      if (typeof err === 'object' && err !== null && 'status' in err) {
        const error = err as { status: number | string; data?: unknown };
        if (error.status === 413) {
          errorMsg = 'File is too large. Please upload a smaller file or contact your administrator to increase the upload limit.';
        } else if (typeof error.status === 'number' && error.status >= 500) {
          errorMsg = 'Server error. Please try again later.';
        }
      }
      
      setErrorMessage(errorMsg);
    }
  };

  const handleDelete = async (recordingId: string) => {
    try {
      await deleteRecording(recordingId).unwrap();
      setDeleteDialogOpen(null);
    } catch (err) {
      console.error('Failed to delete recording:', err);
    }
  };

  const handleDownload = async (fileId: string) => {
    setDownloadingFileId(fileId);
    try {
      await downloadFileWithAuth('zoom-recordings', fileId);
    } catch (err) {
      console.error('Failed to download recording:', err);
      const errorMsg = getErrorMessage(err, 'Failed to download recording');
      setErrorMessage(errorMsg);
    } finally {
      setDownloadingFileId(null);
    }
  };

  const formatFileSize = (bytes?: number): string => {
    if (!bytes) return 'Unknown size';
    if (bytes < 1024) return `${bytes} B`;
    if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(2)} KB`;
    if (bytes < 1024 * 1024 * 1024) return `${(bytes / (1024 * 1024)).toFixed(2)} MB`;
    return `${(bytes / (1024 * 1024 * 1024)).toFixed(2)} GB`;
  };

  const formatDate = (dateString?: string): string => {
    if (!dateString) return 'N/A';
    try {
      return new Date(dateString).toLocaleString();
    } catch {
      return dateString;
    }
  };

  return (
    <Box sx={{ mt: 3 }}>
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
        <Typography variant="h6">Recordings</Typography>
        <Button
          variant="outlined"
          startIcon={<UploadFileIcon />}
          onClick={handleUploadClick}
          disabled={uploading || creating}
        >
          Upload Recording
        </Button>
        <input
          ref={fileInputRef}
          type="file"
          accept="video/*,audio/*"
          style={{ display: 'none' }}
          onChange={handleFileChange}
        />
      </Box>

      {(uploading || creating) && (
        <Box sx={{ mb: 2 }}>
          <LinearProgress variant="indeterminate" />
          <Typography variant="caption" color="text.secondary" sx={{ mt: 0.5 }}>
            {uploading ? 'Uploading file...' : 'Creating recording...'}
          </Typography>
        </Box>
      )}

      {isLoading && <Typography variant="body2" color="text.secondary">Loading recordings...</Typography>}

      {!isLoading && recordings.length === 0 && (
        <Typography variant="body2" color="text.secondary" sx={{ textAlign: 'center', py: 2 }}>
          No recordings yet. Click "Upload Recording" to add one.
        </Typography>
      )}

      {!isLoading && recordings.length > 0 && (
        <List>
          {recordings.map((recording) => (
            <ListItem
              key={recording.id}
              sx={{
                border: '1px solid #e0e0e0',
                borderRadius: 1,
                mb: 1,
                display: 'flex',
                justifyContent: 'space-between',
              }}
            >
              <ListItemText
                primary={
                  <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                    <Typography variant="body2" sx={{ fontWeight: 500 }}>
                      Recording
                    </Typography>
                    {recording.recordingType && (
                      <Chip label={recording.recordingType} size="small" />
                    )}
                  </Box>
                }
                secondary={
                  <Box sx={{ display: 'flex', flexDirection: 'column', gap: 0.5, mt: 0.5 }}>
                    {recording.recordingStartTime && (
                      <Typography variant="caption" color="text.secondary">
                        Start: {formatDate(recording.recordingStartTime)}
                      </Typography>
                    )}
                    {recording.recordingEndTime && (
                      <Typography variant="caption" color="text.secondary">
                        End: {formatDate(recording.recordingEndTime)}
                      </Typography>
                    )}
                    <Typography variant="caption" color="text.secondary">
                      Size: {formatFileSize(recording.fileSize)}
                    </Typography>
                  </Box>
                }
              />
              <Box sx={{ display: 'flex', gap: 0.5 }}>
                <IconButton
                  size="small"
                  color="primary"
                  onClick={() => handleDownload(recording.fileId)}
                  title="Download"
                  disabled={downloadingFileId === recording.fileId}
                >
                  {downloadingFileId === recording.fileId ? (
                    <CircularProgress size={20} />
                  ) : (
                    <DownloadIcon />
                  )}
                </IconButton>
                <IconButton
                  size="small"
                  color="error"
                  onClick={() => setDeleteDialogOpen(recording.id)}
                  title="Delete"
                >
                  <DeleteIcon />
                </IconButton>
              </Box>
            </ListItem>
          ))}
        </List>
      )}

      <Dialog
        open={deleteDialogOpen !== null}
        onClose={() => setDeleteDialogOpen(null)}
        aria-labelledby="delete-dialog-title"
      >
        <DialogTitle id="delete-dialog-title">Delete Recording</DialogTitle>
        <DialogContent>
          <Typography>
            Are you sure you want to delete this recording? This action cannot be undone.
          </Typography>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setDeleteDialogOpen(null)}>Cancel</Button>
          <Button
            onClick={() => deleteDialogOpen && handleDelete(deleteDialogOpen)}
            color="error"
            variant="contained"
          >
            Delete
          </Button>
        </DialogActions>
      </Dialog>

      <Snackbar
        open={errorMessage !== null}
        autoHideDuration={6000}
        onClose={() => setErrorMessage(null)}
        anchorOrigin={{ vertical: 'bottom', horizontal: 'center' }}
      >
        <Alert onClose={() => setErrorMessage(null)} severity="error" sx={{ width: '100%' }}>
          {errorMessage}
        </Alert>
      </Snackbar>
    </Box>
  );
}

