ENTRY  := ./main.go
TARGET := ./api.out
CC     := go

all: build

build:
	$(CC) build -o $(TARGET) $(ENTRY)

clean:
	rm $(TARGET)

test:
	go test ./...

run:
	go run $(TARGET)
