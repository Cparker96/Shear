-- Create the activity table if it doesn't exist
CREATE TABLE IF NOT EXISTS public.activity (
    username VARCHAR(255) NOT NULL,
    update_type VARCHAR(50),
    date VARCHAR(10),
    PRIMARY KEY (username)
);

-- Create an index on username for faster lookups (already covered by PRIMARY KEY, but explicit)
CREATE INDEX IF NOT EXISTS idx_activity_username ON public.activity(username);

-- Grant permissions to the database user (will use the user specified in POSTGRES_USER)
-- Note: This assumes the user already has permissions on the database
-- The user is created automatically by PostgreSQL from POSTGRES_USER env var

