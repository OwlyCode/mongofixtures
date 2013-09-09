// The package mongofixtures enables to load quickly and easily data into mongo db.
package mongofixtures

import (
	"labix.org/v2/mgo"
)

// A session holds the mongo session (based on labix.org/v2/mgo).
type Session struct {
	MongoSession *mgo.Session
	DatabaseName string
}

// Builds a new Session with the given host and database name.
func Begin(host string, databaseName string) (Session, error) {
	mongoSession, err := mgo.Dial(host)
	session := Session{MongoSession: mongoSession, DatabaseName: databaseName}
	return session, err
}

// Removes the collection identified by collectionName.
func (l *Session) Clean(collectionName string) error {
	err := l.MongoSession.DB(l.DatabaseName).C(collectionName).DropCollection()

	if err != nil && err.Error() == "ns not found" {
		err = nil
	}

	return nil
}

/*
Adds a bson-marshalled version of each documents passed as an argument.
Usage example :
    session.Push("collection", Document{Title:"hello world"}, Document{Title:"Hi there!"})
*/
func (l *Session) Push(collectionName string, documents ...interface{}) error {
	var err error
	for _, document := range documents {
		err = l.MongoSession.DB(l.DatabaseName).C(collectionName).Insert(document)
	}
	return err
}

// Ends a session. Should be called with defer right after Begin :
//  session := Begin("localhost", "mydatabase")
//  defer session.End()
func (l *Session) End() {
	l.MongoSession.Close()
}
