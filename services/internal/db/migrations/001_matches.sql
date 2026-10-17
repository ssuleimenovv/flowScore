-- Finished matches: the last snapshot of each, as the API sent it, so a
-- match outlives the server that followed it
CREATE TABLE matches (
	id          text PRIMARY KEY,
	source      text NOT NULL,
	competition text NOT NULL,
	kickoff_at  timestamptz NOT NULL,
	finished_at timestamptz NOT NULL DEFAULT now(),
	snapshot    jsonb NOT NULL
);

CREATE INDEX matches_by_competition ON matches (competition, kickoff_at DESC);
