CREATE TABLE IF NOT EXISTS books (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    author varchar(100) NOT NULL,
    title text NOT NULL,
    year integer NOT NULL,
    readtime integer NOT NULL,
    genres text[] NOT NULL,
    pagecount integer NOT NULL,
    rating real NOT NULL,
    languages text[] NOT NULL,
    version integer NOT NULL DEFAULT 1
);