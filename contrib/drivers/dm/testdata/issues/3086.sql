CREATE TABLE issue3086_user
(
    id        INT NOT NULL COMMENT 'User ID',
    passport  VARCHAR(45) NOT NULL COMMENT 'User Passport',
    password  VARCHAR(45) DEFAULT NULL COMMENT 'User Password',
    nickname  VARCHAR(45) DEFAULT NULL COMMENT 'User Nickname',
    create_at DATETIME DEFAULT NULL COMMENT 'Created Time',
    update_at DATETIME DEFAULT NULL COMMENT 'Updated Time',
    PRIMARY KEY (id)
);
