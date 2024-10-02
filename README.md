# Forum

## Objectives

The objective of this project is to create a web forum that facilitates communication between users, allows categorization of posts, enables liking and disliking of posts and comments, and implements a filtering mechanism for posts.

## Technologies Used

- SQLite for data storage
- Cookies for user authentication and session management
- Docker for containerization

## Setup Instructions

1. Clone the repository to your local machine.
2. Ensure you have Docker installed.
3. Build the Docker image using the provided Dockerfile.
4. Run the Docker container.

## Database Management

- Use SQLite for storing data such as users, posts, comments, etc.
- Design and implement the database structure based on the entity-relationship diagram provided.
- Utilize SQL queries (SELECT, CREATE, INSERT) for database operations.

## Authentication

- Users can register with a unique email, username, and password.
- Ensure email uniqueness and encrypt passwords for security (bonus task).
- Implement a login session using cookies with an expiration date.
- Consider using UUID for session management (bonus task).

## Communication

- Only registered users can create posts and comments.
- Posts can be associated with one or more categories chosen by the user.
- Posts and comments are visible to all users.
- Non-registered users can view posts and comments but cannot interact.

## Likes and Dislikes

- Only registered users can like or dislike posts and comments.
- The number of likes and dislikes is visible to all users.

## Filtering

- Implement filtering options for posts based on categories, created posts, and liked posts.
- Categories act as subforums for specific topics.
- Filtering options are accessible to registered users only and are user-specific.

## Docker

- Dockerize the forum project for easy deployment and management.
- Follow Docker basics for containerization as mentioned in the provided documentation.

## How to Run

``` $ go run ./cmd/web```

and click on the link "http://localhost:8080/"

