package main

import (
	
	"strings"

	"github.com/rs/zerolog/log"

	"database/sql"
	"os"
	"github.com/gofiber/fiber/v3"
	_ "github.com/mattn/go-sqlite3"
)

var sysDbConnection *sql.DB
var userDbConnections = make(map[string]*sql.DB)

type DbCommand struct {
    Command string `json:"command" form:"name" query:"name" validate:"required"`
	DbName string `json:"dbname" form:"dbname" query:"dbname"`
}

type DbCommandResponse struct {
	Status string `json:"status"`
	Msg string `json:"msg"`
	RowsAffected int64 `json:"rowsAffected"` //TODO: Make fields optional
}

func MountRoutes(app *fiber.App) {
	app.Get("/", func(c fiber.Ctx) error {
        // Send a string response to the client
        return c.SendString("sqlheavy version " + VERSION)
    })

	app.Post("/command", func(c fiber.Ctx) error {
		if sysDbConnection == nil {
			log.Error().Msg("App is not initilized!")
			return nil
		}

		dbCommand := new(DbCommand)

		if err := c.Bind().Body(dbCommand); err != nil {
			return c.JSON(DbCommandResponse{"Failed", err.Error(), -1})
		}

		q, err := ParseSql(dbCommand.Command)
		if err != nil {
			return c.JSON(DbCommandResponse{"Failed", err.Error(), -1})
		}

		if q.CreateDatabase != nil {
			//TODO: check not already there
			if _, err := sysDbConnection.ExecContext(c.Context(),
			"INSERT INTO user_db VALUES (?, ?, ?)", q.CreateDatabase.DatabaseName, USER_DB_PATH, "mode=rwc&_mutex=full&_journal=WAL"); err != nil {
				return c.JSON(DbCommandResponse{"Failed", err.Error(), -1})
			}
			constr := "file:" +  USER_DB_PATH + q.CreateDatabase.DatabaseName + ".db" + "?mode=rwc&_mutex=full&_journal=WAL"
			db, err := sql.Open("sqlite3", constr)
			if err != nil {
				return c.JSON(DbCommandResponse{"Failed", err.Error(), -1})
			}
			userDbConnections[q.CreateDatabase.DatabaseName] = db

		} else if q.DetachDatabase != nil {
			if _, err := sysDbConnection.ExecContext(c.Context(),
			"DELETE FROM user_db WHERE name=?", q.DetachDatabase.DatabaseName, USER_DB_PATH, "mode=rwc&_mutex=full&_journal=WAL"); err != nil {
				return c.JSON(DbCommandResponse{"Failed", err.Error(), -1})
			}
			userDbConnections[q.DetachDatabase.DatabaseName].Close()

			delete(userDbConnections, q.DetachDatabase.DatabaseName)
		} else if q.DropDatabase != nil {
			 row := sysDbConnection.QueryRowContext(c.Context(),
			"DELETE FROM user_db WHERE name=? RETURNING path", q.DetachDatabase.DatabaseName, USER_DB_PATH, "mode=rwc&_mutex=full&_journal=WAL")

			var path string
			if err := row.Scan(&path); err != nil {
				return c.JSON(DbCommandResponse{"Failed", err.Error(), -1})
			}
			userDbConnections[q.DetachDatabase.DatabaseName].Close()
			if err := os.Remove(path + q.DetachDatabase.DatabaseName + ".db"); err != nil {
				return c.JSON(DbCommandResponse{"Failed", err.Error(), -1})
			}

		} else {
			//TODO: We have
			// ExecContext: No return values
			// QueryContext: Query with return values
			// QueryRowContext: Query with one row as return value
			// Seems we should infer which one should we use for each query
			switch {
			case q.Select != nil:
				res, err := userDbConnections[dbCommand.DbName].QueryContext(c.Context(), "SELECT " + strings.Join(q.Select.Rest, " "))
				if err != nil {
					return c.JSON(DbCommandResponse{"Failed", err.Error(), -1})
				}
				cols, err := res.Columns()
				if err != nil {
					return c.JSON(DbCommandResponse{"Failed", err.Error(), -1})
				}
				colCount := len(cols)
				rowPicker :=  make([]any, colCount)
				rows := make([]string,0)
				destRow := make([]string, colCount)
				for i := range destRow {
					rowPicker[i] = &destRow[i]
				}

				for res.Next() {
					if err := res.Scan(rowPicker...); err != nil {
						log.Fatal().Err(err).Msg("Error reading rows")
					}
					rows = append(rows, strings.Join(destRow,","))
				}
				return c.JSON(DbCommandResponse{"OK", strings.Join(rows, ";"), -1})
			case q.Insert != nil:
				res, err := userDbConnections[dbCommand.DbName].ExecContext(c.Context(), "INSERT " + strings.Join(q.Insert.Rest, " "))
				if err != nil {
					return c.JSON(DbCommandResponse{"Failed", err.Error(), -1})
				}
				rowsAffected, err := res.RowsAffected()
				if err != nil {
					return c.JSON(DbCommandResponse{"Failed", err.Error(), -1})
				}
				return c.JSON(DbCommandResponse{"OK", "", rowsAffected})
			case q.Update != nil:
				res, err := userDbConnections[dbCommand.DbName].ExecContext(c.Context(), "UPDATE " + strings.Join(q.Insert.Rest, " "))
				if err != nil {
					return c.JSON(DbCommandResponse{"Failed", err.Error(), -1})
				}
				rowsAffected, err := res.RowsAffected()
				if err != nil {
					return c.JSON(DbCommandResponse{"Failed", err.Error(), -1})
				}
				return c.JSON(DbCommandResponse{"OK", "", rowsAffected})
			case q.Delete != nil:
				//TODO: Handle returning which has value
				res, err := userDbConnections[dbCommand.DbName].ExecContext(c.Context(), "DELETE " + strings.Join(q.Delete.Rest, " "))
				if err != nil {
					return c.JSON(DbCommandResponse{"Failed", err.Error(), -1})
				}
				rowsAffected, err := res.RowsAffected()
				if err != nil {
					return c.JSON(DbCommandResponse{"Failed", err.Error(), -1})
				}
				return c.JSON(DbCommandResponse{"OK", "", rowsAffected})
				
			default:
				query := strings.Join(q.Other, " ")
				res, err := userDbConnections[dbCommand.DbName].ExecContext(c.Context(), query)
				if err != nil {
					return c.JSON(DbCommandResponse{"Failed", err.Error(), -1})
				}
				rowsAffected, err := res.RowsAffected()
				if err != nil {
					return c.JSON(DbCommandResponse{"Failed", err.Error(), -1})
				}
				return c.JSON(DbCommandResponse{"OK", "", rowsAffected})

			}
		}
        return c.JSON(DbCommandResponse{"OK", "", -1})
    })

}