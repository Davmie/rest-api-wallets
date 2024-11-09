\connect wallets program

CREATE TABLE IF NOT EXISTS wallet
(
    id            SERIAL PRIMARY KEY,
    uid    uuid UNIQUE NOT NULL,
    amount         INT         NOT NULL
);