CREATE TABLE test_enum (
    id int8 NOT NULL,
    status int2 DEFAULT 0 NOT NULL,
    CONSTRAINT test_enum_pk PRIMARY KEY (id)
);