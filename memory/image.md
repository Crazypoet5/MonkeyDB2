## Image.go

### ImageTable

The design of ImageTable is only to recovery images. 

### CopyTable

When you allocate a image, the copytable will be insert a row to descripe which dst will be copied from which src,
and when another thread writes to the one, it will writes the same thing to copies, too.

### Usage

Create a image and use the datablock it returns, read and write it using the thread-safe functions Read, Write.
When you want to reallocate its size, use RellocImage before you change your pointer, e.g. 
```
    PrimaryKey.A.Datablock = p
    ...
    pNew := ReallocImage(p) // This would use some time while other thread may access the block with code:
                            // PrimaryKey.A.datablock.Write(15, []byte{6, 9, 5, 7})
    PrimaryKey.A.Datablock = pNew
```