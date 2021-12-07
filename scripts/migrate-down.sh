migrate -path ./migrations/user -database "mysql://root:@tcp(localhost:3306)/users" -verbose down
migrate -path ./migrations/post -database "mysql://root:@tcp(localhost:3306)/posts" -verbose down
migrate -path ./migrations/comment -database "mysql://root:@tcp(localhost:3306)/comments" -verbose down
