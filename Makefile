DB_FILE=app.db
INIT_DB:SQLite.sql

.PHONY: start
start: run
		db/init

.PHONY: run
run:
		air

.PHONY: db/init
db/init:
		@echo "Initializing database $(DB_FILE) using $(INIT_DB)..."
		sqlite3 $(DB_FILE) ".read $(INIT_DB)"
		@echo "Database initialized successfully."

.PHONY: db/clear
db/clear:
	@echo "Removing database file..."
	rm -f $(DB_FILE)