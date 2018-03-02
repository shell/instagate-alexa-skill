# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=main

all: clean build
aws: clean build_aws zip
aws_release: aws publish_aws clean

build:
	node dev.js
	$(GOBUILD) -o $(BINARY_NAME) -v

build_aws:
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME) -v

zip:
	zip $(BINARY_NAME).zip $(BINARY_NAME)

publish_aws:
	aws lambda update-function-code --function-name "InstaGateAlexaSkill" --zip-file fileb://$(BINARY_NAME).zip

update_aws_config:
	@read -p "Instagram Login:" instaLogin; \
	read -p "Instagram password:" instaPassword; \
	read -p "Instagram ClientId:" instaClientId; \
	read -p "Instagram ClientSecret:" instaClientSecret; \
	read -p "Instagram AccessToken:" instaAccessToken; \
	aws lambda update-function-configuration --function-name "InstaGateAlexaSkill" --environment "Variables={InstaLogin=$$instaLogin,InstaPassword=$$instaPassword,InstaClientId=$$instaClientId,InstaClientSecret=$$instaClientSecret,InstaAccessToken=$$instaAccessToken}"

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME).zip

run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)
