package models

import "database/sql"

const dbSchema = `
CREATE TABLE IF NOT EXISTS PageGroups (
  ID INTEGER PRIMARY KEY AUTOINCREMENT,
  Key CHAR(6) UNIQUE,
  Name VARCHAR(25),
  Proto INTEGER,
  System CHAR(2),
  SkipFragment BOOL
);

CREATE UNIQUE INDEX IF NOT EXISTS GroupKeyIdx ON PageGroups (Key);

CREATE TABLE IF NOT EXISTS Domains (
  GroupID INTERGER PRIMARY KEY,
  Pattern VARCHAR(50),
  FOREIGN KEY(GroupID) REFERENCES PageGroups (ID)
  ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS Paths (
  GroupID INTERGER PRIMARY KEY,
  Pattern VARCHAR(50),
  FOREIGN KEY(GroupID) REFERENCES PageGroups (ID)
  ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS Locations (
  ID INTEGER PRIMARY KEY AUTOINCREMENT,
  URL VARCHAR(255) UNIQUE,
  GroupID INTEGER,
  FOREIGN KEY(GroupID) REFERENCES PageGroups(ID)
  ON DELETE CASCADE
);

CREATE UNIQUE INDEX IF NOT EXISTS LocURLIdx ON Locations (URL);

CREATE TABLE IF NOT EXISTS Challenges (
  ID INTEGER PRIMARY KEY,
  Challenge VARCHAR(64) UNIQUE,
  LocID INTEGER,
  At TIME,
  IP VARCHAR(40),
  UID INTEGER,
  FOREIGN KEY(LocID) REFERENCES Locations (ID)
  ON DELETE CASCADE 
);

CREATE TABLE IF NOT EXISTS Solutions (
  ChallengeID INTEGER PRIMARY KEY,
  Nonce INTEGER,
  Reward FLOAT,
  FOREIGN KEY(ChallengeID) REFERENCES Challenges (ID)
  ON DELETE CASCADE
);
`

func InitDb(db *sql.DB) (err error) {
	prepare := func(q string) (stmt *sql.Stmt) {
		// Pass-through on first error
		if err != nil {
			return nil
		}
		stmt, err = db.Prepare(q)
		return stmt
	}
	_, err = db.Exec(dbSchema)
	stmtAddGroup = prepare("INSERT INTO PageGroups (Key, Name, Proto, System, SkipFragment) VALUES (?, ?, ?, ?, ?)")
	stmtGetGroup = prepare("SELECT * FROM PageGroups WHERE Key = ?")
	stmtAddDomainPattern = prepare("INSERT INTO Domains (GroupID, Pattern) VALUES (?, ?)")
	stmtAddPathPattern = prepare("INSERT INTO Paths (GroupID, Pattern) VALUES (?, ?)")
	stmtGetDomainPatterns = prepare("SELECT Pattern FROM Domains WHERE GroupID = ?")
	stmtGetPathPatterns = prepare("SELECT Pattern FROM Paths WHERE GroupID = ?")
	return
}
