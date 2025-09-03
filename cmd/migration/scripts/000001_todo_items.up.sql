CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;
CREATE TABLE todo_items (
                       id uuid DEFAULT uuid_generate_v4() NOT NULL,
                       created_at timestamp with time zone NOT NULL,
                       updated_at timestamp with time zone NOT NULL,
                       deleted_at timestamp with time zone,
                       due_date timestamp NOT NULL,
                       description TEXT NOT NULL
);