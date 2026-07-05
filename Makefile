.PHONY: mongo-up mongo-down

mongo-up:
	docker compose up -d

mongo-down:
	docker compose down