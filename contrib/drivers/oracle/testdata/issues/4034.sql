CREATE TABLE issue4034 (
    id         NUMBER(10) NOT NULL,
    passport   VARCHAR2(255),
    password   VARCHAR2(255),
    nickname   VARCHAR2(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);
