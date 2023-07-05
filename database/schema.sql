CREATE DATABASE foo_db;

USE foo_db;

CREATE TABLE
    foo(
        foo_uuid BINARY(16) DEFAULT(UUID_TO_BIN(UUID())) NOT NULL,
        msg VARCHAR(64) NOT NULL,
        CONSTRAINT pk_foo PRIMARY KEY(foo_uuid)
    );