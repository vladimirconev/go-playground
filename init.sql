BEGIN;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS public.job_offers (
  id SERIAL PRIMARY KEY,
  uuid uuid UNIQUE DEFAULT uuid_generate_v4(),
  company text,
  email text,
  expiration_date timestamp without time zone,
  link_to_offer text,
  details text,
  phone text,
  salary double precision NOT NULL, 
  created_at timestamp without time zone DEFAULT current_timestamp NOT NULL,
  deleted_at timestamp without time zone,
  updated_at timestamp without time zone DEFAULT current_timestamp NOT NULL
);

COMMIT;
