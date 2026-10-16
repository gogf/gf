CREATE TABLE issue3754 (
    id        INT NOT NULL,
    name      VARCHAR(45) DEFAULT NULL,
    create_at DATETIME(0) DEFAULT NULL,
    update_at DATETIME(0) DEFAULT NULL,
    delete_at DATETIME(0) DEFAULT NULL,
    PRIMARY KEY (id)
);