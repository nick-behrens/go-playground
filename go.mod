module github.com/snapdocs/go-playground

go 1.16

require (
	github.com/aws/aws-sdk-go-v2/config v1.8.0
	github.com/aws/aws-sdk-go-v2/service/sqs v1.9.0
	github.com/gofrs/uuid v4.0.0+incompatible // indirect
	github.com/jackc/pgtype v1.8.1 // indirect
	github.com/rs/zerolog v1.23.0
	github.com/snapdocs/go-common v0.0.0-20210630182422-a6da52a95593 // indirect
	gorm.io/driver/postgres v1.1.0 // indirect
	gorm.io/gorm v1.21.14 // indirect
	internal/awssqs v1.0.0
)

replace internal/awssqs => ./internal/awssqs
