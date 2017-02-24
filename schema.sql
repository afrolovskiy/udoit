BEGIN;

-- See
-- https://core.telegram.org/bots/api#chat
-- https://core.telegram.org/bots/api#user

CREATE TABLE IF NOT EXISTS tasks (
    id serial PRIMARY KEY,
    description text NOT NULL,
    chat_id bigint NOT NULL,
    creator_id int NOT NULL,
    assignee_id int
);

COMMIT;
