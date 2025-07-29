-- Initialize database schema
-- This script is run automatically when the PostgreSQL container starts

-- Create extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Create a default admin user (password: admin123)
-- This should be changed in production
-- Note: The actual user creation will be handled by the Go application
-- This script just sets up the database structure

-- Set timezone
SET timezone = 'UTC';
