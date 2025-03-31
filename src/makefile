EXECUTABLE_1 = rogue_game

all: run

run: build_1
	./rogue_game

build_1:
	go build -o $(EXECUTABLE_1)

rebuild: clean build

clean:
	@rm -rf rogue_game *.json
