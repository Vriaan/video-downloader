# build_exec:
# 	go build -o video-downloader cmd/main.go && \
# 	./video-downloader --conf config/video-downloader.yaml ; rm video-downloader \


# Launch the program on desktop, it will run using docker mariadb + mailhog on desktop configuration
# launch_program:
# 	docker-compose run video-downloader make build_exec


test:
	@echo "RUNNING Test with command:\n"
	docker-compose run video-downloader go test -cover ./...
