// nolint: govet
package main

import (
	// "strings"
	// "github.com/alecthomas/kong"
	// "github.com/alecthomas/repr"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	
)

type Boolean bool

func (b *Boolean) Capture(values []string) error {
	*b = values[0] == "TRUE"
	return nil
}

type Sql struct {
	CreateDatabase *CreateDatabase `("CREATE DATABASE" @@`
	DropDatabase *DropDatabase ` | "DROP DATABASE" @@`
	DetachDatabase *DetachDatabase ` | "DETACH DATABASE" @@`
	Select *Select ` | "SELECT" @@`
	Delete *Delete ` | "DELETE" @@`
	Insert *Insert` | "INSERT" @@`
	Update *Update` | "Update" @@`
	Other []string ` | (@Ident | @Number | @String | @Operators )*)`
}

type CreateDatabase struct {
	DatabaseName string `@Ident`
}

type DropDatabase struct {
	DatabaseName string `@Ident`
}

type DetachDatabase struct {
	DatabaseName string `@Ident`
}

type Select struct {
	Rest []string ` (@Ident | @Number | @String | @Operators )*`
}

type Delete struct {
	Rest []string ` (@Ident | @Number | @String | @Operators )*`
}

type Insert struct {
	Rest []string ` (@Ident | @Number | @String | @Operators )*`
}

type Update struct {
	Rest []string ` (@Ident | @Number | @String | @Operators )*`
}

var (
	cli struct {
		SQL string `arg:"" required:"" help:"SQL to parse."`
	}
	sqlLexer = lexer.MustSimple([]lexer.SimpleRule{
		{`Keyword`, `(?i)\b(CREATE DATABASE|DROP DATABASE|DETACH DATABASE|SELECT|INSERT|INTO|DELETE|UPDATE)\b`},
		{`Ident`, `[a-zA-Z_][a-zA-Z0-9_]*`},
		{`Number`, `[-+]?\d*\.?\d+([eE][-+]?\d+)?`},
		{`String`, `'[^']*'|"[^"]*"`},
		{`Operators`, `<>|!=|<=|>=|[-+*/%,.()=<>;]`},
		{"whitespace", `\s+`},
	})
	parser = participle.MustBuild[Sql](
		participle.Lexer(sqlLexer),
	//	participle.Unquote("String"),
		participle.CaseInsensitive("Keyword"),
	)
)

func ParseSql(q string) (v *Sql, err error) {
	return parser.ParseString("", q)

}

// func main() {
// 	ctx := kong.Parse(&cli)
// 	sql, err := parser.ParseString("", cli.SQL)
// 	repr.Println(sql, repr.Indent("  "), repr.OmitEmpty(true))
// 	repr.Println(strings.Join(sql.Other, " "))
// 	ctx.FatalIfErrorf(err)
// }
