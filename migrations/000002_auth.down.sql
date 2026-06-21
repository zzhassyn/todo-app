ALTER TABLE todoapp.users
    DROP CONSTRAINT users_email_format,
    DROP CONSTRAINT users_email_unique;

ALTER TABLE todoapp.users
    DROP COLUMN password_hash,
    DROP COLUMN email;
