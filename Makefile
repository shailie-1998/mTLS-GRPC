client:
	go run client/main.go 

server:
	go run server/main.go 

cert:
	cd cert; ./gen.sh; cd ..