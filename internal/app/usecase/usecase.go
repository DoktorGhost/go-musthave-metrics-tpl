package usecase

import "github.com/DoktorGhost/go-musthave-metrics-tpl/internal/app/storage"

type UsecaseMemStorage struct {
	storage storage.RepositoryDB
}

func NewUsecaseMemStorage(storage storage.RepositoryDB) *UsecaseMemStorage {
	return &UsecaseMemStorage{storage: storage}
}

func (uc *UsecaseMemStorage) UsecaseUpdateGuage(nameMetric string, value float64) {
	uc.storage.UpdateGauage(nameMetric, value)
}

func (uc *UsecaseMemStorage) UsecaseUpdateCounter(nameMetric string, value int64) {
	uc.storage.UpdateCounter(nameMetric, value)
}

func (uc *UsecaseMemStorage) UsecaseRead(nameType, nameMetric string) interface{} {
	return uc.storage.Read(nameType, nameMetric)
}

func (uc *UsecaseMemStorage) UsecaseReadAll() map[string]interface{} {
	return uc.storage.ReadAll()
}
