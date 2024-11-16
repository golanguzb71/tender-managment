CREATE TABLE users
(
    id         SERIAL PRIMARY KEY,
    username   VARCHAR(100)                                         NOT NULL,
    password   VARCHAR(255)                                         NOT NULL,
    email      VARCHAR(150)                                         NOT NULL UNIQUE,
    role       VARCHAR(10) CHECK (role IN ('client', 'contractor')) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE tenders
(
    id              SERIAL PRIMARY KEY,
    client_id       INT          NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    title           VARCHAR(255) NOT NULL,
    description     TEXT,
    deadline        DATE         NOT NULL,
    budget          NUMERIC(15, 2) CHECK (budget > 0),
    status          VARCHAR(10) CHECK (status IN ('open', 'closed', 'awarded')) DEFAULT 'open',
    attachment_path varchar,
    created_at      TIMESTAMP                                                   DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP                                                   DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE bids
(
    id            SERIAL PRIMARY KEY,
    tender_id     INT NOT NULL REFERENCES tenders (id) ON DELETE CASCADE,
    contractor_id INT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    price         NUMERIC(15, 2) CHECK (price > 0),
    delivery_time INT CHECK (delivery_time > 0),
    comments      TEXT,
    status        VARCHAR(10) CHECK (status IN ('submitted', 'rejected')) DEFAULT 'submitted',
    created_at    TIMESTAMP                                               DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP                                               DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE notifications
(
    id          SERIAL PRIMARY KEY,
    user_id     INT  NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    message     TEXT NOT NULL,
    relation_id INT,
    type        VARCHAR(15),
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
