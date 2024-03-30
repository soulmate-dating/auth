CREATE TABLE users
(
    id         uuid,
    email      TEXT UNIQUE NOT NULL,
    password   TEXT,
--     logged_in  BOOLEAN,
    PRIMARY KEY (id)
);