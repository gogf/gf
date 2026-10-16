-- Note: DM has no ON UPDATE CURRENT_TIMESTAMP, updated_at only has the default.
CREATE TABLE issue4034 (
    id INT PRIMARY KEY AUTO_INCREMENT,
    passport VARCHAR(255),
    password VARCHAR(255),
    nickname VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);