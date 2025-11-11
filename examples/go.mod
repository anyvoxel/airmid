module github.com/anyvoxel/airmid/examples

go 1.25.3

require (
	github.com/anyvoxel/airmid/anvil v0.0.0-20251110024612-4694a8aee8b4
	github.com/anyvoxel/airmid/app v0.0.0-20251110024612-4694a8aee8b4
	github.com/anyvoxel/airmid/ioc v0.0.0-20251110024612-4694a8aee8b4
	go.opentelemetry.io/otel v1.38.0
	go.opentelemetry.io/otel/exporters/stdout/stdoutmetric v1.28.0
	go.opentelemetry.io/otel/sdk v1.38.0
	go.opentelemetry.io/otel/sdk/metric v1.38.0
)

require (
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/panjf2000/ants/v2 v2.11.3 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/otel/metric v1.38.0 // indirect
	go.opentelemetry.io/otel/trace v1.38.0 // indirect
	go.uber.org/automaxprocs v1.6.0 // indirect
	golang.org/x/sync v0.16.0 // indirect
	golang.org/x/sys v0.35.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace (
// github.com/anyvoxel/airmid/anvil => ../anvil
// github.com/anyvoxel/airmid/app => ../app
// github.com/anyvoxel/airmid/ioc => ../ioc
)
