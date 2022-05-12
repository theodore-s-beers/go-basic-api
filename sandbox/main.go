package main

import (
	"fmt"
	"log"

	"github.com/theodore-s-beers/go-basic-api/internal/database"
)

func main() {
	c := database.NewClient("db.json")
	err := c.EnsureDB()
	if err != nil {
		log.Fatal(err)
	}

	user, err := c.CreateUser("test@example.com", "password", "john doe", 18)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("user created", user)

	updatedUser, err := c.UpdateUser("test@example.com", "password", "john doe", 18)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("user updated", updatedUser)

	gotUser, err := c.GetUser("test@example.com")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("user gotten", gotUser)

	err = c.DeleteUser("test@example.com")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("user deleted")

	_, err = c.GetUser("test@example.com")
	if err == nil {
		log.Fatal("shouldn't be able to get user that was deleted")
	}
	fmt.Println("user confirmed deleted")

	user, err = c.CreateUser("test@example.com", "password", "john doe", 18)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("user recreated", user)

	post, err := c.CreatePost("test@example.com", "my cat is way too fat")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("post created", post)

	secondPost, err := c.CreatePost("test@example.com", "my cat is getting skinny now")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("second post created", secondPost)

	posts, err := c.GetPosts("test@example.com")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("posts gotten", posts)

	err = c.DeletePost(post.ID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("first post deleted", posts)

	posts, err = c.GetPosts("test@example.com")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("posts gotten", posts)

	err = c.DeletePost(secondPost.ID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("second post deleted", posts)

	posts, err = c.GetPosts("test@example.com")
	if err == nil {
		log.Fatal("shouldn't be able to get posts when user has none")
	}
	fmt.Println("user has no remaining posts")

	err = c.DeleteUser("test@example.com")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("user re-deleted")
}
