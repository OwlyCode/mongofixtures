// The package mongofixtures enables to load quickly and easily data into mongo db.
package mongofixtures

import (
	"github.com/kylelemons/go-gypsy/yaml"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"regexp"
	"strconv"
)

// A session holds the mongo session (based on labix.org/v2/mgo).
type Session struct {
	MongoSession *mgo.Session
	DatabaseName string
	ObjectIds    map[string]bson.ObjectId
}

// Builds a new Session with the given host and database name.
func Begin(host string, databaseName string) (Session, error) {
	mongoSession, err := mgo.Dial(host)
	session := Session{MongoSession: mongoSession, DatabaseName: databaseName}
	session.ObjectIds = make(map[string]bson.ObjectId, 0)
	return session, err
}

// Removes the collection identified by collectionName.
func (l *Session) Clean(collectionName string) error {
	_, err := l.MongoSession.DB(l.DatabaseName).C(collectionName).RemoveAll(nil)

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

func (l *Session) ImportYamlFile(path string) {
	file, err := yaml.ReadFile(path)

	if err != nil {
		panic(err)
	}

	nodes := file.Root.(yaml.Map)

	for collectionName, collectionMap := range nodes {
		cMap := collectionMap.(yaml.Map)
		for _, data := range cMap {
			l.Push(collectionName, l.importNode(data))
		}
	}

}

func (l *Session) importNode(value interface{}) interface{} {

	var result interface{}

	switch value.(type) {
	case yaml.Map:
		result = l.importMap(value.(yaml.Map))
	case yaml.List:
		result = l.importList(value.(yaml.List))
	case yaml.Scalar:
		result = l.importScalar(value.(yaml.Scalar))
	}

	return result
}

func (l *Session) importMap(value yaml.Map) interface{} {
	m := make(map[string]interface{}, 0)
	for key, subvalue := range value {
		m[key] = l.importNode(subvalue)
	}

	return m
}

func (l *Session) importList(value yaml.List) interface{} {
	list := make([]interface{}, 0)
	for _, subvalue := range value {
		list = append(list, l.importNode(subvalue))
	}
	return list
}

func (l *Session) importScalar(value yaml.Scalar) interface{} {
	isHolder, _ := regexp.Match("^__(.*)__$", []byte(value))

	var result interface{}

	if isHolder {
		result = l.getObjectId(string(value))
	} else {

		resBool, isBool := strconv.ParseBool(string(value))
		resFloat, isFloat := strconv.ParseFloat(string(value), 32)
		resInt, isInt := strconv.ParseInt(string(value), 10, 32)

		if isInt == nil {
			return resInt
		}
		if isFloat == nil {
			return resFloat
		}
		if isBool == nil {
			return resBool
		}

		return value
	}

	return result
}

func (l *Session) ImportYamlString(yml string) {

}

func (l *Session) getObjectId(placeholder string) bson.ObjectId {
	value, exists := l.ObjectIds[placeholder]
	if !exists {
		value = bson.NewObjectId()
		l.ObjectIds[placeholder] = value
	}

	return value
}

// Ends a session. Should be called with defer right after Begin :
//  session := Begin("localhost", "mydatabase")
//  defer session.End()
func (l *Session) End() {
	l.MongoSession.Close()
}
