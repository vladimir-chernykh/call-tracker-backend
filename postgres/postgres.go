package postgres

import (
	_ "github.com/lib/pq"
	"database/sql"
	log "github.com/sirupsen/logrus"
	"github.com/vladimir-chernykh/call-tracker-backend/calltracker"
)

type CallService struct {
	DB *sql.DB
}

func (s *CallService) Save(c *calltracker.Call) (*int64, error) {
	tx, err := s.DB.Begin()
	defer func() {
		err := tx.Commit()
		if err != nil {
			log.Error("tx.Commit(): ", err)
			panic(err)
		}
	}()

	if err != nil {
		log.Error("s.DB.Begin(): ", err)
		return nil, err
	}

	pres, err := tx.Query(`
INSERT INTO phones (number, created_at, updated_at)
VALUES ($1, NOW(), NOW())
ON CONFLICT (number) DO UPDATE SET updated_at = NOW()
RETURNING id;
`, c.Phone.Number)
	if err != nil {
		log.Error("tx.Query() phones: ", err)
		return nil, err
	}
	var pid int64
	for pres.Next() {
		if err = pres.Scan(&pid); err != nil {
			log.Error("res.Scan(&pid): ", err)
			return nil, err
		}
	}

	cres, err := tx.Query(`
INSERT INTO calls (phone_id, record, created_at, updated_at)
VALUES ($1, $2, NOW(), NOW())
RETURNING id;
`, pid, c.Audio.Buffer)
	if err != nil {
		log.Error("tx.Query() calls: ", err)
		return nil, err
	}
	for cres.Next() {
		if err = cres.Scan(&c.Id); err != nil {
			log.Error("res.Scan(&pid): ", err)
			return nil, err
		}
	}

	return &c.Id, nil
}
