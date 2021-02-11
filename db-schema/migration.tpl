
-- +migrate Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE SCHEMA IF NOT EXISTS iam;
CREATE TABLE iam.users (
	id BIGSERIAL PRIMARY KEY,
	uuid UUID DEFAULT uuid_generate_v4(),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	deleted_at TIMESTAMP WITH TIME ZONE NULL,
    email VARCHAR NOT NULL DEFAULT '',
	is_banned BOOLEAN NOT NULL DEFAULT FALSE,
    price NUMERIC NOT NULL DEFAULT 0,
    note TEXT NOT NULL DEFAULT '',
    plan_id BIGINT NOT NULL DEFAULT 0
);

-- +migrate Down
DROP TABLE iam.users;

-- boolean
-- smallserial
-- serial
-- bigserial
-- smallint
-- integer
-- bigint
-- numeric
-- decimal
-- varchar
-- text
-- timestamptz
-- bytea