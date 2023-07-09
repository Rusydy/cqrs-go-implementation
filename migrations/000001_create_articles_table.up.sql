-- +migrate Up

CREATE TABLE articles (
    id SERIAL PRIMARY KEY,
    author TEXT NOT NULL,
    title TEXT NOT NULL,
    body TEXT NOT NULL,
    created TIMESTAMP NOT NULL
);