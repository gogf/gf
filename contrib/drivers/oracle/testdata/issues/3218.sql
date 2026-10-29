CREATE TABLE issue3218_sys_config (
    id         NUMBER(10)    NOT NULL,
    name       VARCHAR2(255) DEFAULT NULL,
    value      CLOB,
    created_at TIMESTAMP     DEFAULT NULL,
    updated_at TIMESTAMP     DEFAULT NULL,
    PRIMARY KEY (id),
    UNIQUE (name)
);
