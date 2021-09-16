# This Makefile is an easy way to run common operations.

default: run

run:
	@go run main.go

build:
	@go build -o tftask main.go

deploy-test:
	GOOS=linux go build -o tftask main.go
	-ssh -Y tf pkill tftask
	scp ./tftask tf:/home/ec2-user/
	ssh -Y tf CLOUDTAMERIO_URL=${CLOUDTAMERIO_URL} CLOUDTAMERIO_APIKEY=${CLOUDTAMERIO_APIKEY} TERRAFORM_APIKEY=${TERRAFORM_APIKEY} ./tftask