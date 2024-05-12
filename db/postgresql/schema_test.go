package postgresql

import "embed"

//go:embed files/*.sql
var db embed.FS
