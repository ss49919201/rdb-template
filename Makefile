.PHONY: dc-up
dc-up:
	docker compose up -d

.PHONY: exec-sh
exec-sh:dc-up
	docker exec -it rdb /bin/bash -c "mysql -u root -ppassword"

.PHONEY: init-db
init-db:dc-up
	docker exec -i rdb sh /opt/mysql/setup.sh
