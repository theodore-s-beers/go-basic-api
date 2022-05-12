package database

import (
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/google/uuid"
)

// Export structs

type Client struct {
	Path string
}

type Post struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UserEmail string    `json:"userEmail"`
	Text      string    `json:"text"`
}

type User struct {
	CreatedAt time.Time `json:"createdAt"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Name      string    `json:"name"`
	Age       int       `json:"age"`
}

// Private structs

type databaseSchema struct {
	Users map[string]User `json:"users"`
	Posts map[string]Post `json:"posts"`
}

// Export methods

func (c Client) CreatePost(userEmail, text string) (Post, error) {
	db, err := c.readDB()

	if err != nil {
		return Post{}, err
	}

	_, ok := db.Users[userEmail]

	if !ok {
		return Post{}, errors.New("user doesn't exist")
	}

	post := Post{
		ID:        uuid.New().String(),
		CreatedAt: time.Now().UTC(),
		UserEmail: userEmail,
		Text:      text,
	}

	db.Posts[post.ID] = post

	err = c.updateDB(db)

	return post, err
}

func (c Client) CreateUser(email, password, name string, age int) (User, error) {
	db, err := c.readDB()

	if err != nil {
		return User{}, err
	}

	newUser := User{
		CreatedAt: time.Now().UTC(),
		Email:     email,
		Password:  password,
		Name:      name,
		Age:       age,
	}

	_, ok := db.Users[email]

	if ok {
		return User{}, errors.New("user already exists")
	}

	db.Users[email] = newUser

	err = c.updateDB(db)

	return newUser, err
}

func (c Client) DeletePost(id string) error {
	db, err := c.readDB()

	if err != nil {
		return err
	}

	_, ok := db.Posts[id]

	if !ok {
		return errors.New("post doesn't exist")
	}

	delete(db.Posts, id)

	err = c.updateDB(db)

	return err
}

func (c Client) DeleteUser(email string) error {
	db, err := c.readDB()

	if err != nil {
		return err
	}

	_, ok := db.Users[email]

	if !ok {
		return errors.New("user doesn't exist")
	}

	delete(db.Users, email)

	err = c.updateDB(db)

	return err
}

func (c Client) EnsureDB() error {
	_, err := os.ReadFile(c.Path)

	if err != nil {
		err = c.createDB()
	}

	return err
}

func (c Client) GetPosts(userEmail string) ([]Post, error) {
	db, err := c.readDB()

	if err != nil {
		return []Post{}, err
	}

	_, ok := db.Users[userEmail]

	if !ok {
		return []Post{}, errors.New("user doesn't exist")
	}

	userPosts := make([]Post, 0)

	for _, post := range db.Posts {
		if post.UserEmail == userEmail {
			userPosts = append(userPosts, post)
		}
	}

	if len(userPosts) == 0 {
		return []Post{}, errors.New("user has no posts")
	}

	return userPosts, nil
}

func (c Client) GetUser(email string) (User, error) {
	db, err := c.readDB()

	if err != nil {
		return User{}, err
	}

	user, ok := db.Users[email]

	if !ok {
		return User{}, errors.New("user doesn't exist")
	}

	return user, nil
}

func (c Client) UpdateUser(email, password, name string, age int) (User, error) {
	db, err := c.readDB()

	if err != nil {
		return User{}, err
	}

	user, ok := db.Users[email]

	if !ok {
		return User{}, errors.New("user doesn't exist")
	}

	user = User{
		CreatedAt: user.CreatedAt,
		Email:     email,
		Password:  password,
		Name:      name,
		Age:       age,
	}

	db.Users[email] = user

	err = c.updateDB(db)

	return user, err
}

// Private methods

func (c Client) createDB() error {
	db := databaseSchema{
		Users: make(map[string]User),
		Posts: make(map[string]Post),
	}

	data, err := json.Marshal(db)

	if err != nil {
		return err
	}

	err = os.WriteFile(c.Path, data, 0666)

	return err
}

func (c Client) readDB() (databaseSchema, error) {
	var db databaseSchema

	file, err := os.ReadFile(c.Path)

	if err != nil {
		return databaseSchema{}, err
	}

	err = json.Unmarshal(file, &db)

	if err != nil {
		return databaseSchema{}, err
	}

	return db, nil
}

func (c Client) updateDB(db databaseSchema) error {
	data, err := json.Marshal(db)

	if err != nil {
		return err
	}

	err = os.WriteFile(c.Path, data, 0666)

	return err
}

// Export functions

func NewClient(path string) Client {
	return Client{
		Path: path,
	}
}
