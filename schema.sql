BEGIN;

-- See
-- https://core.telegram.org/bots/api#chat
-- https://core.telegram.org/bots/api#user

CREATE TABLE IF NOT EXISTS tasks (
    id serial PRIMARY KEY,
    id_in_chat int NOT NULL,
    chat_id bigint NOT NULL,
    creator_id int NOT NULL,
    assignee_id int
    description text NOT NULL,
    UNIQUE(chat_id, id_in_chat)
);

COMMIT;
