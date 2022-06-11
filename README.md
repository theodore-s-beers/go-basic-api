# go-basic-api

This is a project that I completed as part of the Golang curriculum on
[boot.dev](https://boot.dev/). The idea was to build a toy "social media
backend," with support for users and posts.

When run, the Go program starts a server and listens on `localhost` port 8080.
It creates a database JSON file (if one is not already present), which will
store user and post data.

A basic [CRUD](https://en.wikipedia.org/wiki/Create,_read,_update_and_delete)
API is provided, with two endpoints: `/users/` and `/posts/`.

The `/users/` endpoint will handle an HTTP request of type POST (to create a
user), GET (to fetch data for a given user), PUT (to update a user), or DELETE
(to delete a user). Users are indexed by email address.

The `/posts/` endpoint will handle an HTTP request of type POST (to create a
post for a given user), GET (to retrieve all posts for a given user), or DELETE
(to delete a post, given its UUID). All of these methods require that the user
already exist.

## Example requests

For instance, we could create a user by making a POST to
<http://localhost:8080/users/> with the following JSON in the request body:

```json
{
  "email": "theo@fake.domain",
  "password": "secret",
  "name": "Theo",
  "age": 99
}
```

We could then retrieve that user data by making a GET request to
<http://localhost:8080/users/theo@fake.domain>.

To add a post for our user, we could make a POST to
<http://localhost:8080/posts/> with the following JSON in the request body:

```json
{
  "userEmail": "theo@fake.domain",
  "text": "Lorem ipsum dolor sit"
}
```

(Again, it will be verified that the user exists.) Then we could retrieve all
posts for this user by making a GET request to
<http://localhost:8080/posts/theo@fake.domain>.

And so forth...
