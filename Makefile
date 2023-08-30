# use: make test-local iter=2 verbose=1  # verbose output
# use: make test-local iter=2            # concise output
# use: make test-local iter=[0-9]+       # all tests
agent-bin = ./cmd/agent/agent
server-bin = ./cmd/server/server
test-local:
	go build -o $(agent-bin) ./cmd/agent
	go build -o $(server-bin) ./cmd/server
	metricstest $(if $(verbose), -test.v) -test.run=^TestIteration$(iter)$$ -agent-binary-path=$(agent-bin) -binary-path=$(server-bin) -server-port=8080 -source-path=. -file-storage-path="/tmp/metrics-db.json"
	rm $(agent-bin) $(server-bin)