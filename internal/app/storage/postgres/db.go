package postgres

import (
	"database/sql"
	"github.com/DoktorGhost/go-musthave-metrics-tpl/internal/app/models"
	"log"
	"sync"
)

type PostgresStorage struct {
	db *sql.DB
	mu sync.RWMutex
}

// NewPostgresRepository создает новый экземпляр PostgresRepository.
func NewPostgresStorage(dsn string) (*PostgresStorage, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	// Создание таблицы, если она не существует
	createTableQuery := `
    CREATE TABLE IF NOT EXISTS metrics (
        id SERIAL PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        mtype VARCHAR(255) NOT NULL CHECK (mtype IN ('gauge', 'counter')),
        delta BIGINT, 
        val DOUBLE PRECISION
    );
    `
	if _, err = db.Exec(createTableQuery); err != nil {
		return nil, err
	}

	return &PostgresStorage{db: db}, nil
}

func (r *PostgresStorage) Read(nameType, nameMetric string) interface{} {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result interface{}
	var err error

	switch nameType {
	case "gauge":
		var val float64
		err = r.db.QueryRow("SELECT val FROM metrics WHERE name = $1", nameMetric).Scan(&val)
		result = val
	case "counter":
		var delta int64
		err = r.db.QueryRow("SELECT delta FROM metrics WHERE name = $1", nameMetric).Scan(&delta)
		result = delta
	default:
		return nil
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		log.Println(err)
		return nil
	}

	return result
}

func (r *PostgresStorage) ReadAll() map[string]interface{} {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make(map[string]interface{})

	query := "SELECT name, mtype, val, delta FROM metrics"
	rows, err := r.db.Query(query)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var metr models.Metrics
		err := rows.Scan(&metr.ID, &metr.MType, &metr.Value, &metr.Delta)
		if err != nil {
			log.Println(err)
			return nil
		}
		if metr.MType == "gauge" {
			result[metr.ID] = metr.Value
		} else if metr.MType == "counter" {
			result[metr.ID] = metr.Delta
		}
	}

	// Проверка на ошибки после цикла
	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil
	}
	return result
}

func (r *PostgresStorage) UpdateGauge(nameMetric string, value float64) {
	res := r.Read("gauge", nameMetric)
	r.mu.Lock()
	defer r.mu.Unlock()

	if res == nil {
		query := "INSERT INTO metrics (name, mtype, val) VALUES ($1, $2, $3)"
		_, err := r.db.Exec(query, nameMetric, "gauge", value)
		if err != nil {
			log.Println(err)
			return
		}
	} else {
		query := "UPDATE metrics SET val = $1 WHERE name = $2"
		_, err := r.db.Exec(query, value, nameMetric)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func (r *PostgresStorage) UpdateCounter(nameMetric string, delta int64) {
	res := r.Read("counter", nameMetric)
	r.mu.Lock()
	defer r.mu.Unlock()

	if res == nil {
		query := "INSERT INTO metrics (name, mtype, delta) VALUES ($1, $2, $3)"
		_, err := r.db.Exec(query, nameMetric, "counter", delta)
		if err != nil {
			log.Println(err)
			return
		}
	} else {
		deltaRes := delta + res.(int64)
		query := "UPDATE metrics SET delta = $1 WHERE name = $2"
		_, err := r.db.Exec(query, deltaRes, nameMetric)
		if err != nil {
			log.Println(err)
			return
		}
	}
}
