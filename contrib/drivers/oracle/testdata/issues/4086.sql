CREATE TABLE issue4086 (
    proxy_id      NUMBER(19) NOT NULL,
    recommend_ids CLOB DEFAULT NULL,
    photos        CLOB DEFAULT NULL,
    PRIMARY KEY (proxy_id)
);
INSERT INTO issue4086 (proxy_id, recommend_ids, photos) VALUES (1, '[584, 585]', 'null');
INSERT INTO issue4086 (proxy_id, recommend_ids, photos) VALUES (2, '[]', NULL);
