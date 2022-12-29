CREATE TABLE IF NOT EXISTS smurls (
	small_url varchar NOT NULL,
	created_at timestamptz NOT NULL,
	modified_at timestamptz NOT NULL,
	long_url varchar NOT NULL,
	admin_url varchar NOT NULL,
	count integer,
	ip_info text[]
	);