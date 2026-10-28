CREATE TABLE images (
    id     SERIAL PRIMARY KEY,
    key    TEXT,
    format TEXT,
    width  INT,
    height INT,
    size   INT,
    make     TEXT,
    model    TEXT,
    taken_at TIMESTAMP,
    lat      DOUBLE PRECISION,
    lng      DOUBLE PRECISION
);
