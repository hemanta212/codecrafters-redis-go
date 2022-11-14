package main

type KeyValueStore struct {
	data map[string]string
}

func NewKeyValueStore() *KeyValueStore {
	return &KeyValueStore{data: map[string]string{}}
}

func (store *KeyValueStore) Get(key string) (string, bool) {
	value, found := store.data[key]
	return value, found
}

func (store *KeyValueStore) Set(key, value string) bool {
	store.data[key] = value
	return true
}
