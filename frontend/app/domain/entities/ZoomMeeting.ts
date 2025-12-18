export interface ZoomMeeting {
  id: string;
  sectionModuleId: string;
  zoomMeetingId: string;
  topic: string;
  startTime?: string;
  duration?: number;
  joinUrl: string;
  startUrl: string;
  password?: string;
  createdAt: string;
  updatedAt: string;
}

