package main

import (
	
	"github.com/rs/zerolog/log"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func InitHeavy() {
	// Connect to sysdb and populate in-mem structures
	sysDbConnectionL, err := sql.Open("sqlite3", SYS_DB_CONNECTION_STR)
	if err != nil {
		log.Fatal().Err(err).Msg("Error opening sysdb")
	}
	sysDbConnection = sysDbConnectionL
	log.Debug().Msg("Opening sysdb at " + SYS_DB_PATH)
	rows, err := sysDbConnection.Query("SELECT * FROM user_db")
	if err != nil {
		log.Fatal().Err(err).Msg("Error querying sysdb")
	}
	if rows == nil {
		log.Debug().Msg("Initializing sysdb")
		_, err := sysDbConnection.Exec("CREATE TABLE user_db (name TEXT, path TEXT, options TEXT)")
		if err != nil {
			log.Fatal().Err(err).Msg("Error initializing sysdb")
		}
		log.Debug().Msg("Created user_db table")

	} else {
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
				userDbConnections[dbname] = db
			}
		}
	}
}