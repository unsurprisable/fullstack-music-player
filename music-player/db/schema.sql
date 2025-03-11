-- table storing every individual song entry
CREATE TABLE IF NOT EXISTS songs (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    -- source TEXT CHECK (source IN ('local', 'spotify')) NOT NULL,
    -- filename TEXT,
    -- spotify_uri TEXT,
    filename TEXT NOT NULL,
    title TEXT,
    artist TEXT,
    album TEXT,
    uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- table storing playlists
CREATE TABLE IF NOT EXISTS playlists (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- join table linking playlists with their associated songs
CREATE TABLE IF NOT EXISTS playlists_songs (
    playlist_id INT REFERENCES playlists(id) ON DELETE CASCADE,
    song_id INT REFERENCES songs(id) ON DELETE CASCADE,
    PRIMARY KEY (playlist_id, song_id)
);