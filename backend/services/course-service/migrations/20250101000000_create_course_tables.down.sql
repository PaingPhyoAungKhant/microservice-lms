-- Drop tables in reverse order of dependencies
DROP TABLE IF EXISTS section_module;
DROP TABLE IF EXISTS course_section;
DROP TABLE IF EXISTS course_offering_instructor;
DROP TABLE IF EXISTS course_offering;
DROP TABLE IF EXISTS course_category;
DROP TABLE IF EXISTS course;
DROP TABLE IF EXISTS category;

-- Drop ENUM types
DROP TYPE IF EXISTS content_status;
DROP TYPE IF EXISTS content_type;
DROP TYPE IF EXISTS section_status;
DROP TYPE IF EXISTS offering_status;
DROP TYPE IF EXISTS offering_type;

