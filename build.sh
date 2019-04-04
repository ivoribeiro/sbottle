for GOOS in darwin linux windows android; do
    for GOARCH in 386 amd64 android; do
        go build -v -o bin/sbottle-$GOOS-$GOARCH
    done
done
