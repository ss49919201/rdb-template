.PHONY: dc-up
dc-up:
	docker compose up -d

.PHONY: exec
exec:dc-up
	docker exec -it rdb /bin/bash

.PHONY: exec-mysql
exec-mysql:dc-up
	docker exec -it rdb /bin/bash -c "mysql -u root -ppassword rdb"

.PHONY: which
which:dc-up
	docker exec -i rdb sh /opt/mysql/which.sh

.PHONEY: init-db
init-db:dc-up
	docker exec -i rdb sh /opt/mysql/setup.sh
