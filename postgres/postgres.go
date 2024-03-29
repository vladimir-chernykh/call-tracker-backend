package postgres

import (
	_ "github.com/lib/pq"
	"database/sql"
	log "github.com/sirupsen/logrus"
	"github.com/vladimir-chernykh/call-tracker-backend/calltracker"
	"strconv"
	"time"
	"io/ioutil"
	"fmt"
)

type Storage struct {
	DB *sql.DB
}

func New(db *sql.DB) calltracker.CallStorage {
	return &Storage{DB: db}
}

func (s *Storage) Save(c *calltracker.Call) (*int64, error) {
	log.Info("Save", c.Phone)
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

func (s *Storage) Dump(c *calltracker.Call) (*string, error) {
	log.Info("Save", c.Phone, c.Id)
	aacFilename := strconv.FormatInt(time.Now().UnixNano(), 10) + ".aac"

	rows, err := s.DB.Query("SELECT record FROM calls WHERE id = $1", c.Id)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer rows.Close()
	var record []byte
	rows.Next()
	if err := rows.Scan(&record); err != nil {
		panic(err)
	}
	if err := rows.Err(); err != nil {
		panic(err)
	}
	wErr := ioutil.WriteFile(aacFilename, record, 0644)
	if wErr != nil {
		panic(wErr)
	}

	return &aacFilename, nil
}

func (s *Storage) SaveMetric(m *calltracker.Metric) (error) {
	log.Info("SaveMetric", m.Call.Phone, m.Call.Id)
	tx, err := s.DB.Begin()
	if err != nil {
		panic(err)
	}

	defer func() {
		err := tx.Commit()
		if err != nil {
			panic(err)
		}
	}()

	rs, qErr := tx.Query(`
INSERT INTO metrics (name, call, data, remote_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, NOW(), NOW())
ON CONFLICT (name, call) DO UPDATE SET updated_at = NOW(), data = $3
RETURNING id;
`, m.Name, m.Call.Id, m.Data, m.Call.RemoteId)
	if qErr != nil {
		panic(qErr)
	}
	rs.Close()

	return nil
}

func (s *Storage) GetMetrics(callId string) ([]byte, error) {
	log.Info("GetMetrics ", callId)

	rows, err := s.DB.Query(`
SELECT jsonb_object_agg(name, data :: JSON -> 'result') AS metrics
FROM metrics
WHERE call = $1
GROUP BY call;
`, callId)
	if err != nil {
		log.Error(err)
		panic(err)
	}
	defer rows.Close()

	var record []byte
	if !rows.Next() {
		return nil, nil
	}
	if err := rows.Scan(&record); err != nil {
		log.Error(err)
		panic(err)
	}
	if err := rows.Err(); err != nil {
		log.Error(err)
		panic(err)
	}

	return record, nil
}

func (s *Storage) GetMetricNames() ([]string, error) {
	log.Info("GetMetricNames")

	rows, err := s.DB.Query(`SELECT name FROM metric_names;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []string{}
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}

		out = append(out, name)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}
