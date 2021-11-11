# This Makefile is an easy way to run common operations.

default: run

run:
	@go run main.go

build:
	@go build -o tftask main.go

# Deploy to AWS using CloudFormation.
deploy:
	aws cloudformation create-stack --stack-name tftask --template-body file://cloudformation/ec2-deploy.yaml \
	--parameters file://cloudformation/parameters.json \
	--capabilities CAPABILITY_IAM

# Test requires an SSH profile named 'tf' to test deploying to an AWS EC2 instance.
test:
	GOOS=linux go build -o tftask main.go
	-ssh -Y tf pkill tftask
	scp ./tftask tf:/home/ec2-user/
	ssh -Y tf KION_URL=${KION_URL} KION_APIKEY=${KION_APIKEY} TERRAFORM_APIKEY=${TERRAFORM_APIKEY} ./tftask