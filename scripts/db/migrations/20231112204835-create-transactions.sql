-- +migrate Up
CREATE SEQUENCE public.transactions_id_seq AS bigint;
CREATE TABLE public.transactions (
  id bigint DEFAULT nextval('public.transactions_id_seq') NOT NULL,
  account_id bigint NOT NULL,
  operation_type_id int NOT NULL,
  amount numeric(20,2) NOT NULL,
  event_date timestamp with time zone
);

ALTER TABLE public.transactions
  ADD CONSTRAINT transactions_pkey PRIMARY KEY (id);
ALTER TABLE public.transactions
  ADD CONSTRAINT transactions_account_id_fkey FOREIGN KEY (account_id)
  REFERENCES public.accounts(id) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE public.transactions
  ADD CONSTRAINT transactions_operation_type_id_fkey FOREIGN KEY (operation_type_id)
  REFERENCES public.operation_types(id) ON UPDATE CASCADE ON DELETE CASCADE;
CREATE INDEX transactions_account_idx
  ON public.transactions USING btree (account_id, event_date);

-- +migrate Down
DROP TABLE public.transactions;
DROP SEQUENCE public.transactions_id_seq;
