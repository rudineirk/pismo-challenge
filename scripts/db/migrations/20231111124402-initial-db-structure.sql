-- +migrate Up
CREATE TYPE public.operation_types_enum AS ENUM (
  'CASH_PURCHASE',
  'INSTALLMENT_PURCHASE',
  'WITHDRAWAL',
  'PAYMENT'
);

CREATE SEQUENCE public.accounts_id_seq AS bigint;
CREATE TABLE public.accounts (
  id bigint DEFAULT nextval('public.accounts_id_seq') NOT NULL,
  document_number character varying(255) NOT NULL,
  created_at timestamp with time zone NOT NULL,
  updated_at timestamp with time zone NOT NULL
);

ALTER TABLE public.accounts
  ADD CONSTRAINT accounts_pkey PRIMARY KEY (id);

CREATE INDEX accounts_document_number_idx
  ON public.accounts USING btree (document_number);

CREATE SEQUENCE public.transactions_id_seq AS bigint;
CREATE TABLE public.transactions (
  id bigint DEFAULT nextval('public.transactions_id_seq') NOT NULL,
  account_id bigint NOT NULL,
  operation_type public.operation_types_enum,
  amount numeric(20,2) NOT NULL,
  event_date timestamp with time zone
);

ALTER TABLE public.transactions
  ADD CONSTRAINT transactions_pkey PRIMARY KEY (id);
ALTER TABLE public.transactions
  ADD CONSTRAINT transactions_account_id_fkey FOREIGN KEY (account_id)
  REFERENCES public.accounts(id) ON UPDATE CASCADE ON DELETE CASCADE;
CREATE INDEX transactions_account_idx
  ON public.transactions USING btree (account_id, event_date);

-- +migrate Down
DROP TABLE public.transactions;
DROP TABLE public.accounts;

DROP SEQUENCE public.transactions_id_seq;
DROP SEQUENCE public.accounts_id_seq;

DROP TYPE public.operation_types_enum;
