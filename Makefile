setup:
	export GOOGLE_APPLICATION_CREDENTIALS=C:\Users\uvednara\Documents\playground\gcpCreds\gcpstep-cred.json
clean:
	rm -rf src\.serverless

deploy: clean build
	serverless deploy -c src\serverless.yml

remove: clean
	serverless remove -c src\serverless.yml
	rm -rf src\.serverless
