package database

import (
	"errors"
	"time"
)

type User struct {
	ID                 int       `json:"id"`
	Email              string    `json:"email"`
	Password           string    `json:"password"`
	RefreshToken       string    `json:"refresh_token"`
	RefreshTokenExpiry time.Time `json:"refresh_token_expiry"`
	IsChirpyRed        bool      `json:"is_chirpy_red"`
}

// CreateUser creates a new user and saves it to disk
func (db *DB) CreateUser(email string, password string) (User, error) {
	database, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	if database.Users == nil {
		database.Users = make(map[int]User)
	}
	id := len(database.Users) + 1

	user := User{
		ID:          id,
		Email:       email,
		Password:    password,
		IsChirpyRed: false,
	}
	database.Users[id] = user
	err = db.writeDB(database)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (db *DB) GetUser(id int) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	user, ok := dbStructure.Users[id]
	if !ok {
		return User{}, errors.New("user not found")
	}
	return user, nil
}

// GetUsers returns all users in the database
func (db *DB) GetUsers() ([]User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err

	}
	users := make([]User, 0, len(dbStructure.Users))
	for _, user := range dbStructure.Users {
		users = append(users, user)
	}
	return users, nil
}
func (db *DB) GetUserByEmail(email string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	for _, user := range dbStructure.Users {
		if user.Email == email {
			return user, nil
		}
	}
	return User{}, errors.New("user not found")
}

func (db *DB) UpdateUser(id int, email string, password string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	user, ok := dbStructure.Users[id]
	if !ok {
		return User{}, errors.New("user not found")
	}
	user.Email = email
	user.Password = password
	dbStructure.Users[id] = user

	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (db *DB) UpgradeToChirpyRed(id int) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}
	user, ok := dbStructure.Users[id]
	if !ok {
		return errors.New("user not found")
	}
	user.IsChirpyRed = true
	dbStructure.Users[id] = user
	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) StoreRefreshToken(id int, refreshToken string, expiry time.Time) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}
	user, ok := dbStructure.Users[id]
	if !ok {
		return errors.New("user not found")

	}
	user.RefreshToken = refreshToken
	user.RefreshTokenExpiry = expiry

	dbStructure.Users[id] = user
	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) ValidateRefreshToken(refreshToken string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	for _, user := range dbStructure.Users {
		if user.RefreshToken == refreshToken {
			if time.Now().After(user.RefreshTokenExpiry) {
				return User{}, errors.New("refresh token has expired")
			}
			return user, nil
		}
	}
	return User{}, errors.New("user not found")
}

func (db *DB) DeleteRefreshToken(id int) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}
	user, ok := dbStructure.Users[id]
	if !ok {
		return errors.New("user not found")
	}
	user.RefreshToken = ""
	user.RefreshTokenExpiry = time.Time{}
	dbStructure.Users[id] = user
	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}
	return nil
}
