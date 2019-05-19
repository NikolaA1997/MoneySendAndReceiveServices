package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

type Storage struct {
	Balance   float64   `json:"balance"`
	UpdatedAt time.Time `json:"updatedAt"`
	filepath  string
	mutex     sync.Mutex
}

func (storage Storage) ConvertBalance() float64 {
	balance := storage.Balance / 100
	return balance
}

func InitStorage(storageFile *string) (*Storage, error) {
	storage := &Storage{
		filepath: *storageFile,
	}
	_, err := os.Create(storage.filepath)
	if err != nil {
		return nil, err
	}
	err = storage.Set(0, time.Time{})
	if err != nil {
		return nil, err
	}
	return storage, nil
}

func (s *Storage) load() (*Storage, error) {

	storageJson, _ := json.Marshal(s)
	err := ioutil.WriteFile(s.filepath, storageJson, 0644)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Storage) Set(amount float64, time time.Time) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if time.IsZero() {
		s.Balance = amount
		s.UpdatedAt = time
	} else {
		s.Balance += amount
		s.UpdatedAt = time
	}
	_, err := s.load()
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) Get() float64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	data, _ := ioutil.ReadFile(s.filepath)
	json.Unmarshal(data, s)
	return s.ConvertBalance()
}
