-- +migrate Up
CREATE TABLE public.operation_types (
  id int NOT NULL,
  description character varying(255) NOT NULL
);
ALTER TABLE public.operation_types
  ADD CONSTRAINT operation_types_pkey PRIMARY KEY (id);

INSERT INTO public.operation_types (id, description) VALUES
  (1, 'CASH PURCHASE'),
  (2, 'INSTALLMENT PURCHASE'),
  (3, 'WITHDRAWAL'),
  (4, 'PAYMENT');

-- +migrate Down
DROP TYPE public.operation_types;
