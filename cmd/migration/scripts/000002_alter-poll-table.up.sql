CREATE TABLE public.polls (
    id uuid DEFAULT uuid_generate_v4() NOT NULL,
    version bigint DEFAULT 0 NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    deleted_at timestamp with time zone,
    title  varchar(500) NOT NULL,
    astat varchar(50) NOT NULL
);

CREATE TABLE public.poll_options (
                              id uuid DEFAULT uuid_generate_v4() NOT NULL,
                              version bigint DEFAULT 0 NOT NULL,
                              created_at timestamp with time zone NOT NULL,
                              updated_at timestamp with time zone NOT NULL,
                              deleted_at timestamp with time zone,
                              astat varchar(50) NOT NULL,
                              poll_id  uuid NOT NULL,
                              content text NOT NULL,
                              counts bigint default 0 not null
);