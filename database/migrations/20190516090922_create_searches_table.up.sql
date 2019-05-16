CREATE TABLE searches (
    topics           TEXT DEFAULT NULL,
    languages        TEXT DEFAULT NULL,
    response_status  INT(3) DEFAULT NULL,
    response_content TEXT DEFAULT NULL,
    created_at       DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
