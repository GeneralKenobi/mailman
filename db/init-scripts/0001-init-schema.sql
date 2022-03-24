CREATE SCHEMA mailmandb;
SET search_path TO mailmandb;

CREATE TABLE customer
(
    id    SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL CHECK (email <> ''),

    CONSTRAINT unique_email UNIQUE (email)
);
CREATE INDEX customer_email ON customer (email);

CREATE TABLE mailing_entry
(
    id          SERIAL PRIMARY KEY,
    customer_id INT          NOT NULL,
    mailing_id  INT          NOT NULL,
    title       VARCHAR(255) NOT NULL CHECK (title <> ''),
    content     TEXT,
    insert_time TIMESTAMP    NOT NULL,

    CONSTRAINT fk_customer FOREIGN KEY (customer_id) REFERENCES customer (id)
);
CREATE INDEX mailing_entry_insert_time ON mailing_entry (insert_time);
