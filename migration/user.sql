CREATE TABLE user_demo (
	id smallserial NOT NULL,
	username varchar(30) NOT NULL,
	email varchar(60) NULL,
	"password" bytea NULL,
	"token" varchar NULL,
	role varchar not null ,
	CONSTRAINT admin_pkey PRIMARY KEY (id),
	CONSTRAINT admin_token_key UNIQUE (token),
	CONSTRAINT admin_username_key UNIQUE (username)
);