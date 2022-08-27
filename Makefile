.PHONY: dc-up
dc-up:
	docker compose up -d

.PHONY: exec-sh
exec-sh:dc-up
	docker exec -it rdb /bin/bash -c "mysql -u user -ppassword"
