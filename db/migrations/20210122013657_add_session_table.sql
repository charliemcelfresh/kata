-- migrate:up
CREATE TABLE session
(
  uuid uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  value TEXT NOT NULL,
  user_id uuid
);

-- migrate:down
DROP TABLE session
