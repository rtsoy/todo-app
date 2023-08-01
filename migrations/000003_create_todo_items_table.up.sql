CREATE TABLE todo_items
(
    id          UUID         NOT NULL PRIMARY KEY,
    title       VARCHAR(255) NOT NULL,
    description VARCHAR(255) NOT NULL,
    created_at  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deadline    TIMESTAMP,
    completed   BOOLEAN      NOT NULL DEFAULT FALSE
);

CREATE TABLE lists_items
(
    id      UUID                                              NOT NULL PRIMARY KEY,
    list_id UUID REFERENCES todo_lists (id) ON DELETE CASCADE NOT NULL,
    item_id UUID REFERENCES todo_items (id) ON DELETE CASCADE NOT NULL
);

