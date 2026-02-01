ALTER TABLE users ADD COLUMN gender VARCHAR(10);
ALTER TABLE users ADD CONSTRAINT check_gender CHECK (gender IN ('male', 'female', 'other'));
