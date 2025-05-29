CREATE TABLE public.votes (
    id uuid DEFAULT uuid_generate_v4() NOT NULL,
    version bigint DEFAULT 0 NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    deleted_at timestamp with time zone,
    user_id  uuid NOT NULL,
    poll_id  uuid NOT NULL,
    option_id  uuid NOT NULL,
    astat varchar(50) NOT NULL
);