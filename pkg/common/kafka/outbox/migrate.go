package outbox

import (
	"fmt"

	"github.com/thealiakbari/hichapp/pkg/common/db"
	"gorm.io/gorm"
)

const checkUpCode = `
SELECT 
	COUNT(*) 
FROM "information_schema"."tables" 
WHERE 
	"table_schema" = '%s' 
AND 
	"table_name" = 'outbox_messages'`

const outbox_sql_code = `
CREATE TABLE %s.outbox_messages (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    aggregate_id uuid not null,
    aggregate_type character varying not null,
    trace_id uuid not null,
    type character varying(255) not null,
    name character varying not null,
    payload jsonb not null,
    created_at timestamp without time zone DEFAULT now() NOT NULL
);

ALTER TABLE ONLY %s."outbox_messages"
	 ADD CONSTRAINT "PK_OUTBOX_MESSAGES" PRIMARY KEY (id);
`

func setup(db *gorm.DB, schema string) error {
	err := db.Exec(fmt.Sprintf(outbox_sql_code, schema, schema)).Error
	return err
}

func checkUp(db *gorm.DB) error {
	var schemaName string
	err := db.Raw(`SELECT CURRENT_SCHEMA()`).Scan(&schemaName).Error
	if err != nil {
		return err
	}

	var count int64
	err = db.Raw(fmt.Sprintf(checkUpCode, schemaName)).Scan(&count).Error
	if err != nil {
		return err
	}

	if count == 0 {
		err = setup(db, schemaName)
		if err != nil {
			return err
		}
	}
	return nil
}

func OutboxCheckup(db db.DBWrapper) error {
	err := checkUp(db.DB)
	if err != nil {
		return err
	}
	return nil
}
