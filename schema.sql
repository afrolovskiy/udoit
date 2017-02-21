BEGIN;

CREATE TABLE IF NOT EXISTS tasks (
    id serial PRIMARY KEY,
    description text NOT NULL,
    creator_id bigint NOT NULL,
    assignee_id bigint
);

COMMIT;
