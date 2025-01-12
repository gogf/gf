create table `%s`(
    id         INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
    passport   VARCHAR(45)  NOT NULL DEFAULT passport,
    password   VARCHAR(128) NOT NULL DEFAULT password,
    nickname   VARCHAR(45),
    created_at TIMESTAMP,
    updated_at TIMESTAMP
)