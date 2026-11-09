CREATE TABLE issue3754 (
    id        NUMBER(10)   NOT NULL,
    name      VARCHAR2(45) DEFAULT NULL,
    create_at TIMESTAMP(0) DEFAULT NULL,
    update_at TIMESTAMP(0) DEFAULT NULL,
    delete_at TIMESTAMP(0) DEFAULT NULL,
    PRIMARY KEY (id)
);
