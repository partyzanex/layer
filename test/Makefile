.PHONY: test
test:
	docker-compose -f $(CURDIR)/docker-compose-postgres.yml up -d
	docker-compose -f $(CURDIR)/docker-compose-test.yml up
	docker-compose -f $(CURDIR)/docker-compose-postgres.yml stop
	docker-compose -f $(CURDIR)/docker-compose-postgres.yml rm -f -v
	docker-compose -f $(CURDIR)/docker-compose-test.yml rm -f -v