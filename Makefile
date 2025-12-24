-include .env
export

ENDPOINT    ?= $(PISENSOR_ENDPOINT)
INTERVAL    ?= $(PISENSOR_INTERVAL)
TIMEOUT     ?= $(PISENSOR_TIMEOUT)
SENSOR_ID   ?= $(PISENSOR_SENSOR_ID)
SENSOR_NAME ?= $(PISENSOR_SENSOR_NAME)

BIN_DIR  := bin
BIN      := $(BIN_DIR)/pisensor

LOG_DIR  := logs
LOG      := $(LOG_DIR)/pisensor.log
PID      := $(LOG_DIR)/pisensor.pid

.PHONY: build clean start stop install logs

$(BIN_DIR):
	@mkdir -p $@

clean:
	@echo "[CLEAN] Start"
	@-rm -rf $(BIN_DIR)
	@go clean
	@echo "[CLEAN] Done"

build: clean | $(BIN_DIR)
	@echo "[BUILD] Start"
	@go build -o $(BIN) ./cmd/agent/main.go
	@echo "[BUILD] Done"

start: stop build 
	@echo "[START] Start"
	@mkdir -p "$(LOG_DIR)"
	@nohup $(BIN) -endpoint=$(ENDPOINT) -interval=$(INTERVAL) -timeout=$(TIMEOUT) -sensor_id=$(SENSOR_ID) \
	    >>"$(LOG)" 2>&1 </dev/null & echo $$! >"$(PID)"
	@echo "PID: $$(cat $(PID))  Logs: $(LOG)"
	@echo "[START] Done"

stop: 
	@echo "[STOP] Start"
	@if [ -f "$(PID)" ]; then \
		kill "$$(cat $(PID))" 2>/dev/null || true; \
		rm -f "$(PID)"; \
	else \
		echo "No PID file."; \
	fi
	@echo "[STOP] Done"

install:
	@echo "[INSTALL] Start"
	@bash ./scripts/install.sh
	@echo "[INSTALL] Done"

logs:
	@echo "Tailing $(LOG)… (Ctrl-C to stop)"
	@tail -f "$(LOG)"
