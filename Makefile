clean:
	rm -rf src/.serverless

deploy: clean build
	serverless deploy

remove: clean
	serverless remove
	rm -rf src/.serverless
