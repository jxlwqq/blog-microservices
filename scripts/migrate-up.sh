migrate -path ./migrations/user -database "mysql://root:@tcp(localhost:3306)/users" -verbose up
migrate -path ./migrations/post -database "mysql://root:@tcp(localhost:3306)/posts" -verbose up
migrate -path ./migrations/comment -database "mysql://root:@tcp(localhost:3306)/comments" -verbose up
