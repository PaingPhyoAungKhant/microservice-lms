-- Create ENUM types
CREATE TYPE offering_type AS ENUM ('online', 'oncampus');
CREATE TYPE offering_status AS ENUM ('pending', 'active', 'ongoing', 'completed');
CREATE TYPE section_status AS ENUM ('draft', 'published', 'archived');
CREATE TYPE content_type AS ENUM ('zoom');
CREATE TYPE content_status AS ENUM ('draft', 'pending', 'created');

-- Create category table
CREATE TABLE IF NOT EXISTS category (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    description VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create course table
CREATE TABLE IF NOT EXISTS course (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    thumbnail_id UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create course_category junction table
CREATE TABLE IF NOT EXISTS course_category (
    id UUID PRIMARY KEY,
    course_id UUID NOT NULL REFERENCES course(id) ON DELETE CASCADE,
    category_id UUID NOT NULL REFERENCES category(id) ON DELETE CASCADE,
    UNIQUE(course_id, category_id)
);

-- Create course_offering table
CREATE TABLE IF NOT EXISTS course_offering (
    id UUID PRIMARY KEY,
    course_id UUID NOT NULL REFERENCES course(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    offering_type offering_type NOT NULL,
    status offering_status NOT NULL DEFAULT 'pending',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create course_offering_instructor table
CREATE TABLE IF NOT EXISTS course_offering_instructor (
    id UUID PRIMARY KEY,
    course_offering_id UUID NOT NULL REFERENCES course_offering(id) ON DELETE CASCADE,
    instructor_id UUID NOT NULL,
    instructor_username VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(course_offering_id, instructor_id)
);

-- Create course_section table
CREATE TABLE IF NOT EXISTS course_section (
    id UUID PRIMARY KEY,
    course_offering_id UUID NOT NULL REFERENCES course_offering(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description VARCHAR(255),
    "order" INT NOT NULL DEFAULT 0,
    status section_status NOT NULL DEFAULT 'draft',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create section_module table
CREATE TABLE IF NOT EXISTS section_module (
    id UUID PRIMARY KEY,
    course_section_id UUID NOT NULL REFERENCES course_section(id) ON DELETE CASCADE,
    content_id UUID,
    name VARCHAR(255) NOT NULL,
    description VARCHAR(255),
    content_type content_type NOT NULL,
    content_status content_status NOT NULL DEFAULT 'draft',
    "order" INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create indexes for category
CREATE INDEX IF NOT EXISTS idx_category_name ON category(name);

-- Create indexes for course
CREATE INDEX IF NOT EXISTS idx_course_name ON course(name);
CREATE INDEX IF NOT EXISTS idx_course_thumbnail_id ON course(thumbnail_id);
CREATE INDEX IF NOT EXISTS idx_course_created_at ON course(created_at);

-- Create indexes for course_category
CREATE INDEX IF NOT EXISTS idx_course_category_course_id ON course_category(course_id);
CREATE INDEX IF NOT EXISTS idx_course_category_category_id ON course_category(category_id);

-- Create indexes for course_offering
CREATE INDEX IF NOT EXISTS idx_course_offering_course_id ON course_offering(course_id);
CREATE INDEX IF NOT EXISTS idx_course_offering_status ON course_offering(status);
CREATE INDEX IF NOT EXISTS idx_course_offering_type ON course_offering(offering_type);
CREATE INDEX IF NOT EXISTS idx_course_offering_created_at ON course_offering(created_at);

-- Create indexes for course_offering_instructor
CREATE INDEX IF NOT EXISTS idx_course_offering_instructor_offering_id ON course_offering_instructor(course_offering_id);
CREATE INDEX IF NOT EXISTS idx_course_offering_instructor_instructor_id ON course_offering_instructor(instructor_id);
CREATE INDEX IF NOT EXISTS idx_course_offering_instructor_username ON course_offering_instructor(instructor_username);

-- Create indexes for course_section
CREATE INDEX IF NOT EXISTS idx_course_section_offering_id ON course_section(course_offering_id);
CREATE INDEX IF NOT EXISTS idx_course_section_status ON course_section(status);
CREATE INDEX IF NOT EXISTS idx_course_section_order ON course_section(course_offering_id, "order");

-- Create indexes for section_module
CREATE INDEX IF NOT EXISTS idx_section_module_section_id ON section_module(course_section_id);
CREATE INDEX IF NOT EXISTS idx_section_module_content_id ON section_module(content_id);
CREATE INDEX IF NOT EXISTS idx_section_module_content_status ON section_module(content_status);
CREATE INDEX IF NOT EXISTS idx_section_module_order ON section_module(course_section_id, "order");

