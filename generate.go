package iana

//go:generate go run cmd/iana-name-gen/main.go --url https://www.iana.org/assignments/http-fields/field-names.csv --package field --output field/http.go
//go:generate go run cmd/iana-name-gen/main.go --url https://www.iana.org/assignments/media-types/application.csv --package media --prefix application --output media/application.go
//go:generate go run cmd/iana-name-gen/main.go --url https://www.iana.org/assignments/media-types/audio.csv --package media --prefix audio --output media/audio.go
//go:generate go run cmd/iana-name-gen/main.go --url https://www.iana.org/assignments/media-types/font.csv --package media --prefix font --output media/font.go
//go:generate go run cmd/iana-name-gen/main.go --url https://www.iana.org/assignments/media-types/image.csv --package media --prefix image --output media/image.go
//go:generate go run cmd/iana-name-gen/main.go --url https://www.iana.org/assignments/media-types/message.csv --package media --prefix message --output media/message.go
//go:generate go run cmd/iana-name-gen/main.go --url https://www.iana.org/assignments/media-types/model.csv --package media --prefix model --output media/model.go
//go:generate go run cmd/iana-name-gen/main.go --url https://www.iana.org/assignments/media-types/multipart.csv --package media --prefix multipart --output media/multipart.go
//go:generate go run cmd/iana-name-gen/main.go --url https://www.iana.org/assignments/media-types/text.csv --package media --prefix text --output media/text.go
//go:generate go run cmd/iana-name-gen/main.go --url https://www.iana.org/assignments/media-types/video.csv --package media --prefix video --output media/video.go
//go:generate go fmt ./field
//go:generate go fmt ./media
