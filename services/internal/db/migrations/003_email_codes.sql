-- A user signs in once their email is confirmed by a code. The users who
-- signed up before codes existed count as confirmed
ALTER TABLE users ADD COLUMN email_verified_at timestamptz;
UPDATE users SET email_verified_at = created_at;

-- The codes sent by email, kept as a hash: one live code per user and
-- purpose, and a new one replaces the one before
CREATE TABLE email_codes (
	user_id    uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
	purpose    text NOT NULL CHECK (purpose IN ('verify', 'reset')),
	code_hash  text NOT NULL,
	expires_at timestamptz NOT NULL,
	attempts   int NOT NULL DEFAULT 0,
	sent_at    timestamptz NOT NULL DEFAULT now(),
	PRIMARY KEY (user_id, purpose)
);
