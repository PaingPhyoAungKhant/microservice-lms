INSERT INTO users (id, email, username, password_hash, role, status, email_verified, email_verified_at, created_at, updated_at)
VALUES 
(
    gen_random_uuid(),
    'instructor@asto-lms.local',
    'instructor',
    '$2a$10$XDsC/7TcgBXSLu3ONOQ.wutiaSRWf3/yYyol4yK0hxQbo5wP38nkW',
    'instructor',
    'active',
    true,
    NOW(),
    NOW(),
    NOW()
),
(
    gen_random_uuid(),
    'student@asto-lms.local',
    'student',
    '$2a$10$XDsC/7TcgBXSLu3ONOQ.wutiaSRWf3/yYyol4yK0hxQbo5wP38nkW',
    'student',
    'active',
    true,
    NOW(),
    NOW(),
    NOW()
)
ON CONFLICT (email) DO NOTHING;

