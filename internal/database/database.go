package database

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

type Chirp struct {
	ID       int    `json:"id"`
	Body     string `json:"body"`
	AuthorID int    `json:"author_id"`
}

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users  map[int]User  `json:"users"`
}

// NewDB creates a new database connection
// and creates the database file if it doesn't exist
func NewDB(path string) (*DB, error) {
	db := &DB{
		path: path,
		mux:  &sync.RWMutex{},
	}

	err := db.ensureDB()
	return db, err
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string, author_id int) (Chirp, error) {
	database, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}
	id := len(database.Chirps) + 1

	chirp := Chirp{
		ID:       id,
		Body:     body,
		AuthorID: author_id,
	}
	database.Chirps[id] = chirp
	err = db.writeDB(database)
	if err != nil {
		return Chirp{}, err
	}
	return chirp, nil
}

func (db *DB) DeleteChirp(id int) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}
	_, ok := dbStructure.Chirps[id]
	if !ok {
		return errors.New("chirp not found")
	}
	delete(dbStructure.Chirps, id)
	return db.writeDB(dbStructure)
}

func (db *DB) GetChirp(id int) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}
	chirp, ok := dbStructure.Chirps[id]
	if !ok {
		return Chirp{}, errors.New("chirp not found")
	}
	return chirp, nil
}

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err

	}
	chirps := make([]Chirp, 0, len(dbStructure.Chirps))
	for _, chirp := range dbStructure.Chirps {
		chirps = append(chirps, chirp)
	}
	return chirps, nil
}

// ensureDB creates a new database file if it doesn't exist

func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return db.createDB()
	}
	return err
}

func (db *DB) createDB() error {
	dbStructure := DBStructure{
		Chirps: make(map[int]Chirp),
		Users:  make(map[int]User),
	}
	return db.writeDB(dbStructure)
}

// loadDB reads the database file into memory
func (db *DB) loadDB() (DBStructure, error) {

	db.mux.RLock()
	defer db.mux.RUnlock()

	dbStructure := DBStructure{}
	dat, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return dbStructure, err
	}
	err = json.Unmarshal(dat, &dbStructure)
	if err != nil {
		return dbStructure, err
	}
	return dbStructure, nil
}

// writeDB writes the database file to disk
func (db *DB) writeDB(dbStructure DBStructure) error {
	db.mux.Lock()
	defer db.mux.Unlock()

	dat, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}

	err = os.WriteFile(db.path, dat, 0600)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) ResetDB() error {
	err := os.Remove(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return db.ensureDB()
}
