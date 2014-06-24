package db

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"os"
	"reflect"
	"strings"
	"testing"
)

type MongoDBTestSuite struct {
	suite.Suite
	mongoUrl string
}

func TestRunnerMongoDB(t *testing.T) {
	suite.Run(t, new(MongoDBTestSuite))
}

func (suite *MongoDBTestSuite) SetupSuite() {
	suite.mongoUrl = os.Getenv("MONGO_URL")
}

func (suite *MongoDBTestSuite) TestNew() {
	mongodb := NewMongoDB(suite.mongoUrl)

	assert.Equal(suite.T(), suite.mongoUrl, mongodb.Url())
}

func (suite *MongoDBTestSuite) TestConnection() {
	mongodb := NewMongoDB(suite.mongoUrl)

	err := mongodb.Dial()
	defer mongodb.Close()

	assert.NoError(suite.T(), err)
}

func (suite *MongoDBTestSuite) TestDatabase() {
	dbName := suite.mongoUrl[strings.LastIndex(suite.mongoUrl, "/")+1:]
	mongodb := NewMongoDB(suite.mongoUrl)

	err := mongodb.Dial()
	defer mongodb.Close()

	assert.NoError(suite.T(), err)

	db, err := mongodb.Database()

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), dbName, db.Name)
	assert.Equal(suite.T(), "*mgo.Database", reflect.TypeOf(db).String())
}

func (suite *MongoDBTestSuite) TestDatabaseNameEmptyError() {
	mongodb := NewMongoDB(suite.mongoUrl[:strings.LastIndex(suite.mongoUrl, "/")])

	err := mongodb.Dial()
	defer mongodb.Close()

	_, err = mongodb.Database()

	assert.Error(suite.T(), err)
}

func (suite *MongoDBTestSuite) TestConnectionWithoutPrefix() {
	mongodb := NewMongoDB(suite.mongoUrl[10:])

	err := mongodb.Dial()
	defer mongodb.Close()

	assert.NoError(suite.T(), err)
}
