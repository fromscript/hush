-- Create rooms table
CREATE TABLE IF NOT EXISTS rooms (
                                   id UUID PRIMARY KEY,
                                   created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  last_activity TIMESTAMP NOT NULL DEFAULT NOW()
  );

-- Create messages table
CREATE TABLE IF NOT EXISTS messages (
                                      id UUID PRIMARY KEY,
                                      room_id UUID REFERENCES rooms(id) ON DELETE CASCADE,
  content BYTEA NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  is_read BOOLEAN NOT NULL DEFAULT false
  );

CREATE INDEX IF NOT EXISTS messages_room_id_idx ON messages(room_id);
