package mongofixtures

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"strconv"
	"testing"
)

type document struct {
	Id    bson.ObjectId `bson:"_id,omitempty"`
	Title string        `bson:"title"`
}

func checkCount(t *testing.T, collection *mgo.Collection, search string, count int, message string) {
	c, err := collection.Find(bson.M{"title": search}).Count()

	if err != nil {
		t.Fatal(err)
	}

	if c != count {
		t.Fatal(message + " : " + strconv.Itoa(c))
	}
}

func TestLoader(t *testing.T) {

	mongoSession, err := mgo.Dial("localhost")
	collection := mongoSession.DB("sample").C("collection1")

	session, err := Begin("localhost", "sample")
	defer session.End()

	if err != nil {
		t.Fatal(err)
	}

	err = session.Push("collection1", document{Id: bson.NewObjectId(), Title: "This is a demo"})

	if err != nil {
		t.Fatal(err)
	}

	checkCount(t, collection, "This is a demo", 1, "Wrong count after inserting a document")

	err = session.Push("collection1", document{Id: bson.NewObjectId(), Title: "This is a demo 2"})

	if err != nil {
		t.Fatal(err)
	}

	checkCount(t, collection, "This is a demo 2", 1, "Wrong count after inserting a document")

	err = session.Clean("collection1")

	if err != nil {
		t.Fatal(err)
	}

	checkCount(t, collection, "This is a demo", 0, "Wrong count after freeing the collection")
}
