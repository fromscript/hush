CREATE TABLE rooms (
                     id UUID PRIMARY KEY,
                     encryption_key BYTEA NOT NULL
);

CREATE TABLE messages (
                        id SERIAL PRIMARY KEY,
                        room_id UUID REFERENCES rooms(id),
                        content BYTEA NOT NULL,
                        sender VARCHAR(36),
                        timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                        is_read BOOLEAN DEFAULT FALSE
);
