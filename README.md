# nbt-go

Go interfaces for reading, modifying, and writing Named Binary Tag format data

Read the API:  [nbt.go](https://github.com/mischanix/nbt-go/blob/master/nbt.go)

[GoDoc](http://godoc.org/github.com/Mischanix/nbt-go)

## Example

    // Basic reading:
    testFile, _ := os.Open("bigtest.nbt")
    testStream, _ := gzip.NewReader(bigtestFile)
    test, _ := nbt.Load(bigtestStream)
    test.Name() // "Level"
    test.At("shortTest").Short() // int16(32767)
    test.Path("nested compound test/ham/name").StringData() // "Hampus"
    // ...
    //
    // Back and forth:
    root := nbt.MakeRoot("")
    root.Set("testString", "This is a test String!")
    byteBuf := bytes.NewBuffer(...)
    root.Save(byteBuf)
    byteBuf.Bytes() // [10 0 0 8 0 10 116 101 ...]
    nbt.Load(bytes.NewBuffer(byteBuf.Bytes())) // Same as root
