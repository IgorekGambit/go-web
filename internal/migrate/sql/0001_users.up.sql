CREATE TABLE users (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	session_id TEXT NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	CONSTRAINT users_session_id_key UNIQUE (session_id)
);
