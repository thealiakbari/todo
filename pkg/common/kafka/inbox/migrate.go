package inbox

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
	"table_name" = 'inbox_messages'`

const inbox_sql_code = `
CREATE TABLE %s."inbox_messages" (
  id uuid NOT NULL,
  aggregate_id character varying(256) NOT NULL,
  aggregate_type character varying(256) NOT NULL,
  payload jsonb,
  state character varying(256) DEFAULT 'none'::character varying NOT NULL,
  status character varying(256) DEFAULT 'none'::character varying NOT NULL,
  metadata jsonb,
  created_at timestamp with time zone DEFAULT now() NOT NULL
);
CREATE INDEX idx_inbox_messages_id ON %s.inbox_messages (aggregate_type);
CREATE INDEX idx_inbox_messages_aggregate_type ON %s.inbox_messages (aggregate_type);
`

const missed_columns = `
ALTER TABLE %s.inbox_messages ALTER COLUMN id SET DEFAULT public.uuid_generate_v4();
ALTER TABLE %s.inbox_messages ADD COLUMN IF NOT EXISTS retry_count INT DEFAULT 0 NOT NULL;
ALTER TABLE %s.inbox_messages ADD COLUMN IF NOT EXISTS version INT DEFAULT 1 NOT NULL;
ALTER TABLE %s.inbox_messages ADD COLUMN IF NOT EXISTS wait_duration INT;
ALTER TABLE %s.inbox_messages ADD COLUMN IF NOT EXISTS correlation_id varchar(40) DEFAULT 'OLD_ROWS' NOT NULL;
ALTER TABLE %s.inbox_messages ADD COLUMN IF NOT EXISTS trace_id varchar(40) DEFAULT 'OLD_ROWS' NOT NULL;
`

func setupV1(db *gorm.DB, schema string) error {
	err := db.Exec(fmt.Sprintf(inbox_sql_code, schema, schema, schema)).Error
	return err
}

func setupV2(db *gorm.DB, schema string) error {
	err := db.Exec(fmt.Sprintf(missed_columns, schema, schema, schema, schema, schema, schema)).Error
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
		err = setupV1(db, schemaName)
		if err != nil {
			return err
		}
		err = setupV2(db, schemaName)
		if err != nil {
			return err
		}
	} else {
		err = setupV2(db, schemaName)
		if err != nil {
			return err
		}
	}

	return nil
}

func InboxCheckup(db db.DBWrapper) error {
	err := checkUp(db.DB)
	if err != nil {
		return err
	}
	return nil
}
