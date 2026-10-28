CREATE TABLE issue3086_user (
    id        NUMBER(10)   NOT NULL,
    passport  VARCHAR2(45) NOT NULL,
    password  VARCHAR2(45) DEFAULT NULL,
    nickname  VARCHAR2(45) DEFAULT NULL,
    create_at TIMESTAMP    DEFAULT NULL,
    update_at TIMESTAMP    DEFAULT NULL,
    PRIMARY KEY (id)
);
