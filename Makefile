clean:
	rm -rf ./bin
	rm -rf .serverless

deploy: clean build
	serverless deploy

remove: clean
	serverless remove
	rm -rf .serverless
