CREATE TABLE users
(
    id         TEXT,
    email      TEXT UNIQUE NOT NULL,
    password   TEXT,
--     logged_in  BOOLEAN,
    PRIMARY KEY (id)
);