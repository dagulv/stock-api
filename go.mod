module github.com/dagulv/stock-api

go 1.22.4

require (
	github.com/dagulv/ticker v0.0.1
	github.com/jackc/pgx/v5 v5.6.0
	github.com/json-iterator/go v1.1.12
	github.com/rubenv/sql-migrate v1.7.0
	gopkg.in/yaml.v2 v2.4.0
)

replace github.com/dagulv/ticker v0.0.1 => ../ticker

require (
	github.com/fxamacker/cbor/v2 v2.7.0 // indirect
	github.com/go-gorp/gorp/v3 v3.1.0 // indirect
	github.com/go-webauthn/x v0.1.12 // indirect
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	github.com/golang-jwt/jwt/v5 v5.2.1 // indirect
	github.com/google/go-tpm v0.9.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	github.com/labstack/gommon v0.4.2 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/robfig/cron/v3 v3.0.1 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	golang.org/x/crypto v0.25.0 // indirect
	golang.org/x/net v0.27.0 // indirect
	golang.org/x/sync v0.7.0 // indirect
	golang.org/x/sys v0.22.0 // indirect
	golang.org/x/text v0.16.0 // indirect
	golang.org/x/time v0.5.0 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
)

require (
	github.com/go-webauthn/webauthn v0.10.2
	github.com/joho/godotenv v1.5.1
	github.com/labstack/echo/v4 v4.12.0
	github.com/rs/xid v1.5.0
)
