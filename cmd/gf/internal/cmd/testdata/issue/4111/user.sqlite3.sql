create table `%s`(
    f_id         INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
    f_passport   VARCHAR(45)  NOT NULL DEFAULT passport,
    f_nickname   VARCHAR(45),
    f_created_at TIMESTAMP
)
