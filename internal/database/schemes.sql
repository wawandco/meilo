CREATE TABLE IF NOT EXISTS emails
(
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    subject     TEXT,
    sender      TEXT,
    recipients  TEXT NOT NULL, -- JSON array
    cc          TEXT,          -- JSON array
    bcc         TEXT,          -- JSON array
    received_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS bodies
(
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    email_id     INTEGER NOT NULL REFERENCES emails (id) ON DELETE CASCADE,
    content_type TEXT,
    content      TEXT
);

CREATE TABLE IF NOT EXISTS attachments
(
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    email_id     INTEGER NOT NULL REFERENCES emails (id) ON DELETE CASCADE,
    name         TEXT,
    content_type TEXT,
    data         BLOB
);