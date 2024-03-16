-- run migrate
migrate -database "postgresql://postgres:postgres@localhost:5432/test_mkp?sslmode=disable" -path "./migrate/" up