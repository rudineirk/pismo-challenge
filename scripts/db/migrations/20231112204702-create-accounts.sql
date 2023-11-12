-- +migrate Up
CREATE SEQUENCE public.accounts_id_seq AS bigint;
CREATE TABLE public.accounts (
  id bigint DEFAULT nextval('public.accounts_id_seq') NOT NULL,
  document_number character varying(255) NOT NULL,
  created_at timestamp with time zone NOT NULL,
  updated_at timestamp with time zone NOT NULL
);

ALTER TABLE public.accounts
  ADD CONSTRAINT accounts_pkey PRIMARY KEY (id);

CREATE UNIQUE INDEX accounts_document_number_key
  ON public.accounts USING btree (document_number);

-- +migrate Down
DROP TABLE public.accounts;
DROP SEQUENCE public.accounts_id_seq;
