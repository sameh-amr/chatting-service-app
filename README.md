## 🚀 Getting Started

1. Create a .env file in the project root with the following content:

DB_HOST=db
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=secret
DB_NAME=chatdb

2. Build and Start with Docker Compose:

docker-compose up --build

3. Verify the Setup:

The app should print: Successfully connected to database!

App will be available on: http://localhost:8080

5. Stop the App:

docker-compose down

To also remove volumes:
docker-compose down -v
