import React, { useState, useEffect } from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  TextField,
  Box,
  Typography,
  Divider,
  Alert,
} from '@mui/material';
import DeleteIcon from '@mui/icons-material/Delete';
import EditIcon from '@mui/icons-material/Edit';
import VideoCallIcon from '@mui/icons-material/VideoCall';
import LoginIcon from '@mui/icons-material/Login';
import type { SectionModule } from '../../../domain/entities/SectionModule';
import type { ZoomMeeting } from '../../../domain/entities/ZoomMeeting';
import {
  useGetZoomMeetingByModuleQuery,
  useCreateZoomMeetingMutation,
  useUpdateZoomMeetingMutation,
  useDeleteZoomMeetingMutation,
} from '../../../infrastructure/api/rtk/zoomMeetingsApi';
import ZoomRecordingSection from './ZoomRecordingSection';

interface ModuleContentModalProps {
  open: boolean;
  onClose: () => void;
  module: SectionModule | null;
  mode: 'create' | 'manage';
}

export default function ModuleContentModal({
  open,
  onClose,
  module,
  mode,
}: ModuleContentModalProps) {
  const [topic, setTopic] = useState('');
  const [startDate, setStartDate] = useState<string>('');
  const [startTime, setStartTime] = useState<string>('');
  const [duration, setDuration] = useState<number | ''>('');
  const [password, setPassword] = useState('');
  const [isEditing, setIsEditing] = useState(false);
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);

  const { data: meeting, isLoading: meetingLoading } = useGetZoomMeetingByModuleQuery(
    module?.id || '',
    {
      skip: !module || !open || mode === 'create' || !module.contentId,
    }
  );

  const [createMeeting, { isLoading: creating }] = useCreateZoomMeetingMutation();
  const [updateMeeting, { isLoading: updating }] = useUpdateZoomMeetingMutation();
  const [deleteMeeting, { isLoading: deleting }] = useDeleteZoomMeetingMutation();

  useEffect(() => {
    if (mode === 'create' && module) {
      setTopic(module.name);
      setStartDate('');
      setStartTime('');
      setDuration('');
      setPassword('');
      setIsEditing(false);
    } else if (mode === 'manage' && meeting) {
      setTopic(meeting.topic);
      if (meeting.startTime) {
        const date = new Date(meeting.startTime);
        const localDate = new Date(date.getTime() - date.getTimezoneOffset() * 60000);
        const dateStr = localDate.toISOString().slice(0, 10);
        const timeStr = localDate.toISOString().slice(11, 16);
        setStartDate(dateStr);
        setStartTime(timeStr);
      } else {
        setStartDate('');
        setStartTime('');
      }
      setDuration(meeting.duration || '');
      setPassword(meeting.password || '');
      setIsEditing(false);
    }
  }, [mode, module, meeting, open]);

  const handleCreate = async () => {
    if (!module || !topic.trim()) return;

    try {
      let formattedStartTime: string | null = null;
      if (startDate && startTime) {
        const localDateTime = new Date(`${startDate}T${startTime}`);
        formattedStartTime = localDateTime.toISOString().replace(/\.\d{3}Z$/, 'Z');
      }

      await createMeeting({
        sectionModuleId: module.id,
        topic: topic.trim(),
        startTime: formattedStartTime,
        duration: duration ? Number(duration) : null,
        password: password.trim() || null,
      }).unwrap();
      onClose();
    } catch (err) {
      console.error('Failed to create zoom meeting:', err);
    }
  };

  const handleUpdate = async () => {
    if (!meeting || !topic.trim()) return;

    try {
      let formattedStartTime: string | null = null;
      if (startDate && startTime) {
        const localDateTime = new Date(`${startDate}T${startTime}`);
        formattedStartTime = localDateTime.toISOString().replace(/\.\d{3}Z$/, 'Z');
      }

      await updateMeeting({
        id: meeting.id,
        data: {
          topic: topic.trim(),
          startTime: formattedStartTime,
          duration: duration ? Number(duration) : null,
          password: password.trim() || null,
        },
      }).unwrap();
      setIsEditing(false);
    } catch (err) {
      console.error('Failed to update zoom meeting:', err);
    }
  };

  const handleDelete = async () => {
    if (!meeting) return;

    try {
      await deleteMeeting(meeting.id).unwrap();
      setDeleteDialogOpen(false);
      onClose();
    } catch (err) {
      console.error('Failed to delete zoom meeting:', err);
    }
  };

  const handleStartMeeting = () => {
    if (meeting?.startUrl) {
      window.open(meeting.startUrl, '_blank');
    }
  };

  const handleJoinMeeting = () => {
    if (meeting?.joinUrl) {
      window.open(meeting.joinUrl, '_blank');
    }
  };

  const isLoading = meetingLoading || creating || updating || deleting;
  const isCreateMode = mode === 'create';
  const hasMeeting = !isCreateMode && meeting && !meetingLoading;

  return (
    <Dialog open={open} onClose={onClose} maxWidth="md" fullWidth>
      <DialogTitle>
        {isCreateMode
          ? `Create Zoom Meeting - ${module?.name || ''}`
          : `Manage Zoom Meeting - ${module?.name || ''}`}
      </DialogTitle>
      <DialogContent>
        {isCreateMode ? (
          <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2, mt: 1 }}>
            <TextField
              label="Topic"
              fullWidth
              required
              value={topic}
              onChange={(e) => setTopic(e.target.value)}
              helperText="Meeting topic (defaults to module name)"
            />
            <Box sx={{ display: 'flex', gap: 2 }}>
              <TextField
                label="Start Date (Optional)"
                type="date"
                fullWidth
                value={startDate}
                onChange={(e) => setStartDate(e.target.value)}
                InputLabelProps={{
                  shrink: true,
                }}
              />
              <TextField
                label="Start Time (Optional)"
                type="time"
                fullWidth
                value={startTime}
                onChange={(e) => setStartTime(e.target.value)}
                InputLabelProps={{
                  shrink: true,
                }}
                inputProps={{
                  step: 60,
                }}
              />
            </Box>
            <TextField
              label="Duration (minutes)"
              type="number"
              fullWidth
              value={duration}
              onChange={(e) => setDuration(e.target.value ? Number(e.target.value) : '')}
              inputProps={{ min: 1 }}
              helperText="Optional: Meeting duration in minutes"
            />
            <TextField
              label="Password (Optional)"
              type="password"
              fullWidth
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              helperText="Optional: Meeting password"
            />
          </Box>
        ) : hasMeeting ? (
          <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2, mt: 1 }}>
            {isEditing ? (
              <>
                <TextField
                  label="Topic"
                  fullWidth
                  required
                  value={topic}
                  onChange={(e) => setTopic(e.target.value)}
                />
                <Box sx={{ display: 'flex', gap: 2 }}>
                  <TextField
                    label="Start Date (Optional)"
                    type="date"
                    fullWidth
                    value={startDate}
                    onChange={(e) => setStartDate(e.target.value)}
                    InputLabelProps={{
                      shrink: true,
                    }}
                  />
                  <TextField
                    label="Start Time (Optional)"
                    type="time"
                    fullWidth
                    value={startTime}
                    onChange={(e) => setStartTime(e.target.value)}
                    InputLabelProps={{
                      shrink: true,
                    }}
                    inputProps={{
                      step: 60,
                    }}
                  />
                </Box>
                <TextField
                  label="Duration (minutes)"
                  type="number"
                  fullWidth
                  value={duration}
                  onChange={(e) => setDuration(e.target.value ? Number(e.target.value) : '')}
                  inputProps={{ min: 1 }}
                />
                <TextField
                  label="Password (Optional)"
                  type="password"
                  fullWidth
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                />
              </>
            ) : (
              <>
                <Box sx={{ display: 'flex', gap: 2, mb: 2 }}>
                  <Button
                    variant="contained"
                    startIcon={<VideoCallIcon />}
                    onClick={handleStartMeeting}
                    color="primary"
                  >
                    Start Meeting
                  </Button>
                  <Button
                    variant="outlined"
                    startIcon={<LoginIcon />}
                    onClick={handleJoinMeeting}
                  >
                    Join Meeting
                  </Button>
                </Box>
                <Divider />
                <Box sx={{ display: 'flex', flexDirection: 'column', gap: 1 }}>
                  <Typography variant="body2">
                    <strong>Topic:</strong> {meeting.topic}
                  </Typography>
                  {meeting.startTime && (
                    <Typography variant="body2">
                      <strong>Start Time:</strong>{' '}
                      {new Date(meeting.startTime).toLocaleString()}
                    </Typography>
                  )}
                  {meeting.duration && (
                    <Typography variant="body2">
                      <strong>Duration:</strong> {meeting.duration} minutes
                    </Typography>
                  )}
                  {meeting.password && (
                    <Typography variant="body2">
                      <strong>Password:</strong> {meeting.password}
                    </Typography>
                  )}
                </Box>
                <Divider />
                <ZoomRecordingSection meetingId={meeting.id} />
              </>
            )}
          </Box>
        ) : (
          <Alert severity="info">Loading meeting information...</Alert>
        )}
      </DialogContent>
      <DialogActions>
        {isCreateMode ? (
          <>
            <Button onClick={onClose} disabled={isLoading}>
              Cancel
            </Button>
            <Button
              onClick={handleCreate}
              variant="contained"
              disabled={isLoading || !topic.trim()}
            >
              {creating ? 'Creating...' : 'Create Meeting'}
            </Button>
          </>
        ) : hasMeeting ? (
          isEditing ? (
            <>
              <Button onClick={() => setIsEditing(false)} disabled={isLoading}>
                Cancel
              </Button>
              <Button
                onClick={handleUpdate}
                variant="contained"
                disabled={isLoading || !topic.trim()}
              >
                {updating ? 'Updating...' : 'Save Changes'}
              </Button>
            </>
          ) : (
            <>
              <Button onClick={onClose}>Close</Button>
              <Button
                startIcon={<EditIcon />}
                onClick={() => setIsEditing(true)}
                disabled={isLoading}
              >
                Edit
              </Button>
              <Button
                startIcon={<DeleteIcon />}
                onClick={() => setDeleteDialogOpen(true)}
                color="error"
                disabled={isLoading}
              >
                Delete
              </Button>
            </>
          )
        ) : (
          <Button onClick={onClose}>Close</Button>
        )}
      </DialogActions>

      <Dialog open={deleteDialogOpen} onClose={() => setDeleteDialogOpen(false)}>
        <DialogTitle>Delete Zoom Meeting</DialogTitle>
        <DialogContent>
          <Typography>
            Are you sure you want to delete this zoom meeting? This action cannot be undone.
          </Typography>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setDeleteDialogOpen(false)} disabled={deleting}>
            Cancel
          </Button>
          <Button onClick={handleDelete} color="error" variant="contained" disabled={deleting}>
            {deleting ? 'Deleting...' : 'Delete'}
          </Button>
        </DialogActions>
      </Dialog>
    </Dialog>
  );
}

