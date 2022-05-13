build:
	rm -rf bin
	mkdir bin
	cd src && go build -o ../bin/donggu .
	cp -r templates bin/templates
run: build
	./bin/donggu