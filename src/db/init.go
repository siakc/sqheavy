package db

import (
	"database/sql"
	. "sqheavy/settings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
)

var SysDbConnection *sql.DB
var UserDbConnections = make(map[string]*sql.DB)

func InitHeavy() {
	// Connect to sysdb and populate in-mem structures
	sysDbConnection, err := sql.Open("sqlite3", SYS_DB_CONNECTION_STR)
	if err != nil {
		log.Fatal().Err(err).Msg("Error opening sysdb")
	}
	log.Debug().Msg("Opening sysdb at " + SYS_DB_PATH)
	ensureUserDB(sysDbConnection)
	rows, err := sysDbConnection.Query("SELECT * FROM user_db")
	if err != nil {
		log.Fatal().Err(err).Msg("Error querying sysdb")
	}
	for rows.Next() {
		var dbname string
		var dbpath string
		var options string
		if err := rows.Scan(&dbname, &dbpath, &options); err != nil {
			log.Fatal().Err(err).Msg("Error reading sysdb rows")
		}
		var userdb_conn_str string = "file:" + dbpath + dbname + ".db?" + options
		log.Debug().Msg("Opening user DB: " + userdb_conn_str)
		db, err := sql.Open("sqlite3", userdb_conn_str)
		if err != nil {
			log.Error().Err(err).Msg("Error opening user DB")
		} else {
			UserDbConnections[dbname] = db
		}
	}
}

func ensureUserDB(sysDbConnection *sql.DB) {
	log.Debug().Msg("Initializing sysdb")
	_, err := sysDbConnection.Exec("CREATE TABLE IF NOT EXISTS user_db (name TEXT, path TEXT, options TEXT)")
	if err != nil {
		log.Fatal().Err(err).Msg("Error initializing sysdb")
	}
	log.Debug().Msg("Created user_db table")
}
