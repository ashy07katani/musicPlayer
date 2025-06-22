package repo

const (
	INSERT_SONG = `INSERT INTO songsInfo (title, artist, album, genre, duration,  playlist_path, album_art, created_at ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) ON CONFLICT (title, album) DO NOTHING`
)

/*
 id SERIAL PRIMARY KEY,
musicdb(#     title VARCHAR(255),
musicdb(#     artist VARCHAR(255),
musicdb(#     album VARCHAR(255),
musicdb(#     genre VARCHAR(100),
musicdb(#     duration FLOAT,
musicdb(#     playlist_path TEXT,
musicdb(#     album_art TEXT,
musicdb(#     created_at TIMESTAMP DEFAULT NOW()
*/
