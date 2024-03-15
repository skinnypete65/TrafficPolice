migrateup:
	 migrate -path migrations -database "postgresql://user:user@localhost:5440/traffic_police_db?sslmode=disable" --verbose up

migrateup1:
	 migrate -path migrations -database "postgresql://user:user@localhost:5440/traffic_police_db?sslmode=disable" --verbose up 1

migratedown:
	 migrate -path migrations -database "postgresql://user:user@localhost:5440/traffic_police_db?sslmode=disable" --verbose down

migratedown1:
	 migrate -path migrations -database "postgresql://user:user@localhost:5440/traffic_police_db?sslmode=disable" --verbose down 1

create_cover_out:
	 go test -v -coverprofile cover.out ./...

create_cover_html:
	go tool cover -html cover.out -o cover.html