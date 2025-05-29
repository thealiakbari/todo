CREATE TABLE public.tags (
    id uuid DEFAULT uuid_generate_v4() NOT NULL,
    version bigint DEFAULT 0 NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    deleted_at timestamp with time zone,
    name  varchar(500) NOT NULL,
    astat varchar(50) NOT NULL
);

CREATE TABLE public.poll_tags (
                                  id uuid DEFAULT uuid_generate_v4() NOT NULL,
                                  version bigint DEFAULT 0 NOT NULL,
                                  created_at timestamp with time zone NOT NULL,
                                  updated_at timestamp with time zone NOT NULL,
                                  deleted_at timestamp with time zone,
                                  poll_id  uuid NOT NULL,
                                  tag_id uuid NOT NULL
);