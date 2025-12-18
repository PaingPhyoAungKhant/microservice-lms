export interface ZoomRecording {
  id: string;
  zoomMeetingId: string;
  fileId: string;
  recordingType?: string;
  recordingStartTime?: string;
  recordingEndTime?: string;
  fileSize?: number;
  createdAt: string;
  updatedAt: string;
}

