# nbt-go

Read the API:  [nbt.go](thtps://github.com/mischanix/nbt-go/blob/master/nbt.go)

[GoDoc](http://godoc.org/github.com/mischanix/nbt-go)

## Example

    testFile, _ := os.Open("bigtest.nbt")
    testStream, _ := gzip.NewReader(bigtestFile)
    test, _ := nbt.Load(bigtestStream)
    test.Name() // "Level"
    test.At("shortTest").Short() // int16(32767)
    test.Path("nested compound test/ham/name").StringData() // "Hampus"
    ...
