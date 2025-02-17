BEGIN;

CREATE TABLE IF NOT EXISTS public.users (
	id serial NOT NULL,
	user_name varchar(50) NOT NULL,
	password_hash varchar(255) NOT NULL,
	CONSTRAINT users_pk PRIMARY KEY (id),
	CONSTRAINT users_unique UNIQUE (user_name)
);

CREATE TABLE IF NOT EXISTS public.request (
	id serial NOT NULL,
	contact varchar(255) NOT NULL,
	contact_type varchar(10) NOT NULL,
	message text NOT NULL,
	created_at timestamp with time zone NOT NULL,
	CONSTRAINT request_pk PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS public.outbox (
	request_id serial4 NOT NULL,
	is_sent bool NOT NULL,
	CONSTRAINT outbox_pk PRIMARY KEY (request_id)
);

COMMIT;