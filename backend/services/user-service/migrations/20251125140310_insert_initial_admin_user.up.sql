INSERT INTO users (id, email, username, password_hash, role, status, email_verified, email_verified_at, created_at, updated_at)
VALUES (
    gen_random_uuid(),
    'admin@asto-lms.local',
    'admin',
    '$2a$10$XDsC/7TcgBXSLu3ONOQ.wutiaSRWf3/yYyol4yK0hxQbo5wP38nkW',
    'admin',
    'active',
    true,
    NOW(),
    NOW(),
    NOW()
)
ON CONFLICT (email) DO NOTHING;

