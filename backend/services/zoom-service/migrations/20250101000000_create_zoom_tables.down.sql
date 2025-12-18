DROP INDEX IF EXISTS idx_zoom_recording_created_at;
DROP INDEX IF EXISTS idx_zoom_recording_file_id;
DROP INDEX IF EXISTS idx_zoom_recording_meeting_id;

DROP INDEX IF EXISTS idx_zoom_meeting_created_at;
DROP INDEX IF EXISTS idx_zoom_meeting_zoom_meeting_id;
DROP INDEX IF EXISTS idx_zoom_meeting_section_module_id;

DROP TABLE IF EXISTS zoom_recording;
DROP TABLE IF EXISTS zoom_meeting;

