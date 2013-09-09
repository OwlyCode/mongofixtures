mongofixtures [![Build Status](https://travis-ci.org/OwlyCode/mongofixtures.png)](https://travis-ci.org/OwlyCode/mongofixtures)
===============
A Go quick and dirty utility for cleaning collections and loading fixtures into them.

How to install ?
----------------

```bash
go get github.com/OwlyCode/mongofixtures
```

How to use ?
------------

```go
session, err := Begin("localhost", "database")
defer session.End()

if err != nil {
  panic(err)
}

// Initial cleaning.
session.Clean("collection")

// Now population time !
session.Push("collection", document{Id: bson.NewObjectId(), Title: "This is a demo"})
```

Note that Push and Clean can return both an error that you might want to check.

You can provide to the Begin function host parameter any string representing a mongo connection. For example : "mongodb://myuser:mypass@localhost" should work fine.
