CREATE TABLE todo_lists
(
    id          UUID         NOT NULL PRIMARY KEY,
    title       VARCHAR(255) NOT NULL,
    description VARCHAR(255) NOT NULL,
    created_at  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE users_lists
(
    id      UUID                                              NOT NULL PRIMARY KEY,
    user_id UUID REFERENCES users (id) ON DELETE CASCADE      NOT NULL,
    list_id UUID REFERENCES todo_lists (id) ON DELETE CASCADE NOT NULL
);
