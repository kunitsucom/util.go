package mysql_test

import (
	"log"
	"strings"

	"github.com/kunitsucom/util.go/exp/database/sql/ddl/mysql"
)

func ExampleParser_Parse() {
	p := mysql.NewParser(strings.NewReader(`--   これは先頭行 (DDL 以前) のコメントです。
CREATE TABLE IF NOT EXISTS users -- これは CREATE TABLE 文のコメントです。
(
    -- これは id カラムのこめんとです。
    "id" INTEGER PRIMARY KEY AUTOINCREMENT,
    -- これは name カラムのコメントです。
    "name" TEXT NOT NULL,
    -- これは enable カラムのコメントです。
    "enable" TINYINT(1) NOT NULL DEFAULT 1,
    -- これは age カラムのコメントです。
    "age" INTEGER NULL DEFAULT NULL,
	-- これは created_at カラムのコメントです。
	"created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP(),
	-- これは updated_at カラムのコメントです。
	"updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP() ON UPDATE CURRENT_TIMESTAMP()
)
ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='これはテーブルのコメントです。';
; --これはセミコロンの行のコメントです。
-- これは末尾行 (DDL 以降) のコメントです。`))
	ddl, err := p.Parse()
	if err != nil {
		panic(err)
	}
	log.Printf("%s", ddl)
	// Output:
	//
	log.Printf("%#v", ddl)
}
