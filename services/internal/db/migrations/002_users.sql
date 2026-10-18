-- The people who signed up. The email is kept lower-cased, so the same
-- address typed in another case is the same user
CREATE TABLE users (
	id            uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	email         text NOT NULL UNIQUE,
	name          text NOT NULL,
	password_hash text NOT NULL,
	created_at    timestamptz NOT NULL DEFAULT now()
);


