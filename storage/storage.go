package storage

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

func NewStorage(filepath string) *Storage {
	return &Storage{
		filepath: filepath,
	}
}

func (s *Storage) load() (*Storage, error) {
	f, err := os.Open(s.filepath)
	if err != nil {
		f, _ = os.Create(s.filepath)
		storageJson, _ := json.Marshal(s)
		err = ioutil.WriteFile(s.filepath, storageJson, 0644)
		if err != nil {
			return nil, err
		}
	}
	defer f.Close()

	byteValue, _ := ioutil.ReadAll(f)
	json.Unmarshal([]byte(byteValue), &s)
	return s, nil
}

func (s *Storage) save() error {
	f, err := os.Open(s.filepath)
	if err != nil {
		return err
	}
	defer f.Close()

	storageJson, _ := json.Marshal(s)
	err = ioutil.WriteFile(s.filepath, storageJson, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) Set(amount float64,time time.Time) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	storage, err := s.load()
	if err != nil {
		return err
	}
	if time.IsZero(){}
	{
		storage.Balance = amount
	}
		storage.Balance += amount

	storage.UpdatedAt = time
	err = storage.save()

	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) Get() (float64, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	storage, err := s.load()
	if err != nil {
		return 0, err
	}

	return storage.Balance, nil
}