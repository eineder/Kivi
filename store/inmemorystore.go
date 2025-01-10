package store

import (
	"errors"
	"fmt"
)

type InMemoryStore struct {
	data map[string]string
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		data: make(map[string]string),
	}
}

func (s *InMemoryStore) CreateItem(key string, value string) error {
	if _, exists := s.data[key]; exists {
		return errors.New("item already exists")
	}
	s.data[key] = value
	fmt.Printf("\nCreated item '%s' with value '%s'", key, value)
	return nil
}

func (s *InMemoryStore) GetItem(key string) (string, error) {
	value, exists := s.data[key]
	if !exists {
		return "", errors.New("item not found")
	}
	fmt.Printf("\nGot item '%s' with value '%s'", key, value)
	return value, nil
}

func (s *InMemoryStore) UpdateItem(key string, value string) error {
	if _, exists := s.data[key]; !exists {
		return errors.New("item not found")
	}
	s.data[key] = value
	fmt.Printf("\nUpdated item '%s' with value '%s'", key, value)
	return nil
}

func (s *InMemoryStore) DeleteItem(key string) error {
	if _, exists := s.data[key]; !exists {
		return errors.New("item not found")
	}
	delete(s.data, key)
	fmt.Printf("\nDeleted item '%s'", key)
	return nil
}
