package models

import "database/sql"

const dbSchema = `
CREATE TABLE IF NOT EXISTS locations (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  url VARCHAR(255) UNIQUE NOT NULL,
  hash CHAR(32)
);

CREATE TABLE IF NOT EXISTS solutions (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  lid INTEGER,
  reward FLOAT,
  challenge CHAR(32),
  nonce INTEGER,
  FOREIGN KEY(lid) REFERENCES locations (id)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_loc_url ON locations(url);
`

func InitDb(driver, source string) (db *sql.DB, err error) {
	prepare := func(q string) (stmt *sql.Stmt) {
		// Pass-through on first error
		if err != nil {
			return nil
		}
		stmt, err = db.Prepare(q)
		return stmt
	}
	db, err = sql.Open(driver, source)
	if err != nil {
		return
	}
	_, err = db.Exec(dbSchema)
	stmtGetLoc = prepare("SELECT * FROM locations WHERE url = ?")
	stmtAddLoc = prepare("INSERT INTO locations (url, hash) VALUES (?, ?)")
	return
}
