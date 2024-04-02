SRC_DIR := src
BUILD_DIR := build

build:
	GOOS=linux GOARCH=amd64 go build -tags lambda.norpc -o $(BUILD_DIR)/bootstrap ./$(SRC_DIR)
	cd $(BUILD_DIR) && zip deployment.zip bootstrap

create-function:
	aws lambda create-function --function-name Esep-Webhook \
	--runtime provided.al2 --handler bootstrap \
	--architectures amd64 \
	--role arn:aws:iam::975050140595:role/lambda-ex \
	--zip-file fileb://$(BUILD_DIR)/deployment.zip

update-function:
	aws lambda update-function-code --function-name Esep-Webhook --zip-file fileb://$(BUILD_DIR)/deployment.zip

clean:
	rm -f $(BUILD_DIR)/bootstrap $(BUILD_DIR)/deployment.zip

run:
	go run $(SRC_DIR)/main.go