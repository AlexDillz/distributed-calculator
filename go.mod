module github.com/AlexDillz/distributed-calculator

go 1.23.0

require (
	github.com/golang-jwt/jwt/v5 v5.2.2
	github.com/mattn/go-sqlite3 v1.14.28
	github.com/stretchr/testify v1.10.0
	google.golang.org/grpc v1.72.0
	google.golang.org/protobuf v1.36.5
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/net v0.35.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/text v0.22.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250218202821-56aae31c358a // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	dmitri.shuralyov.com/gpu/mtl => github.com/dmitshur/gpu-mtl v0.0.0-20220304182956-4a5f9b60beed
	github.com/AlexDillz/distributed-calculator/internal/agent => ./internal/agent
	github.com/AlexDillz/distributed-calculator/internal/proto => ./internal/proto
	github.com/AlexDillz/distributed-calculator/internal/server => ./internal/server
	github.com/AlexDillz/distributed-calculator/internal/storage => ./internal/storage
)
