-- migrations/001_init.sql
CREATE TABLE transcripts (
     id INTEGER PRIMARY KEY AUTOINCREMENT,
     title TEXT,
     source_url TEXT,
     backend TEXT,
     language TEXT,
     text TEXT,
     created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
