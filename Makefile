.PHONY:

# ==============================================================================
# Docker

develop:
	echo "Starting develop docker compose"
	docker-compose -f docker-compose.yml up --build -d

local:
	echo "Starting local docker compose"
	docker-compose -f docker-compose.local.yml up --build -d

down:
	echo "Stopping local docker compose"
	docker-compose down --remove-orphans

logs:
	echo "Starting log local docker compose"
	docker-compose logs -f

crate_topics:
	docker exec -it kafka1 kafka-topics --zookeeper zookeeper:2181 --create --topic create-product --partitions 3 --replication-factor 2
	docker exec -it kafka1 kafka-topics --zookeeper zookeeper:2181 --create --topic update-product --partitions 3 --replication-factor 2
	docker exec -it kafka1 kafka-topics --zookeeper zookeeper:2181 --create --topic dead-letter-queue --partitions 3 --replication-factor 2
