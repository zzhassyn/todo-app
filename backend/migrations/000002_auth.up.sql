ALTER TABLE todoapp.users
    ADD COLUMN email VARCHAR(255),
    ADD COLUMN password_hash VARCHAR(255);

UPDATE todoapp.users SET email = 'user' || id || '@todoapp.local' WHERE email IS NULL;
UPDATE todoapp.users SET password_hash = '' WHERE password_hash IS NULL;

ALTER TABLE todoapp.users
    ALTER COLUMN email SET NOT NULL,
    ALTER COLUMN password_hash SET NOT NULL;

ALTER TABLE todoapp.users
    ADD CONSTRAINT users_email_unique UNIQUE (email),
    ADD CONSTRAINT users_email_format CHECK (email ~ '^[^@\s]+@[^@\s]+\.[^@\s]+$');
