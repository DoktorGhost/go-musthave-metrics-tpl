package storage

// Repository представляет интерфейс для работы с хранилищем данных.
type RepositoryDB interface {
	UpdateGauage(nameMetric string, value float64)
	UpdateCounter(nameMetric string, value int64)
	Read(nameType, nameMetric string) interface{}
	ReadAll() map[string]interface{}
}
