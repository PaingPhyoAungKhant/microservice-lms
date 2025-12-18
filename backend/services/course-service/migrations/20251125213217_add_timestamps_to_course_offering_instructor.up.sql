DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'course_offering_instructor' AND column_name = 'created_at') THEN
        ALTER TABLE course_offering_instructor ADD COLUMN created_at TIMESTAMPTZ NOT NULL DEFAULT NOW();
    END IF;
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'course_offering_instructor' AND column_name = 'updated_at') THEN
        ALTER TABLE course_offering_instructor ADD COLUMN updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();
    END IF;
END $$;

