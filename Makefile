build:
	go build -o main .
	go vet
	go fmt
	golint

run:
	go run  .

compile:
	echo "Compiling for every OS and Platform"
	
