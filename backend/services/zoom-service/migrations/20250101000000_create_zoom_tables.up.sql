CREATE TABLE IF NOT EXISTS zoom_meeting (
    id UUID PRIMARY KEY,
    section_module_id UUID NOT NULL,
    zoom_meeting_id VARCHAR(255) NOT NULL UNIQUE,
    topic VARCHAR(255) NOT NULL,
    start_time TIMESTAMPTZ,
    duration INT,
    join_url TEXT NOT NULL,
    start_url TEXT NOT NULL,
    password VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS zoom_recording (
    id UUID PRIMARY KEY,
    zoom_meeting_id UUID NOT NULL REFERENCES zoom_meeting(id) ON DELETE CASCADE,
    file_id UUID NOT NULL,
    recording_type VARCHAR(50),
    recording_start_time TIMESTAMPTZ,
    recording_end_time TIMESTAMPTZ,
    file_size BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_zoom_meeting_section_module_id ON zoom_meeting(section_module_id);
CREATE INDEX IF NOT EXISTS idx_zoom_meeting_zoom_meeting_id ON zoom_meeting(zoom_meeting_id);
CREATE INDEX IF NOT EXISTS idx_zoom_meeting_created_at ON zoom_meeting(created_at);

CREATE INDEX IF NOT EXISTS idx_zoom_recording_meeting_id ON zoom_recording(zoom_meeting_id);
CREATE INDEX IF NOT EXISTS idx_zoom_recording_file_id ON zoom_recording(file_id);
CREATE INDEX IF NOT EXISTS idx_zoom_recording_created_at ON zoom_recording(created_at);

