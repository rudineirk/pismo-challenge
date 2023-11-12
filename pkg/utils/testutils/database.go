package testutils

import (
	"database/sql"
	"fmt"
	"net/url"

	uuid "github.com/satori/go.uuid"
)

type TestDatabase struct {
	URL           string
	managementURL string
	databaseName  string
}

func NewTestDatabase(databaseURL string) (*TestDatabase, error) {
	databaseName := fmt.Sprintf("test-%s", uuid.NewV4().String())
	testdb := TestDatabase{managementURL: databaseURL, databaseName: databaseName}
	testdb.buildURLConnection()

	if err := testdb.createNewDatabase(); err != nil {
		return nil, err
	}

	return &testdb, nil
}

func (testdb *TestDatabase) buildURLConnection() {
	url, err := url.Parse(testdb.managementURL)
	if err != nil {
		panic(err)
	}

	url.Path = testdb.databaseName
	testdb.URL = url.String()
}

func (testdb *TestDatabase) createNewDatabase() error {
	sqlDB, err := sql.Open("postgres", testdb.managementURL)
	if err != nil {
		return err
	}

	defer func() {
		if err := sqlDB.Close(); err != nil {
			return
		}
	}()

	_, err = sqlDB.Exec(fmt.Sprintf("CREATE DATABASE \"%s\"", testdb.databaseName))

	return err
}

func (testdb *TestDatabase) Drop() error {
	sqlDB, err := sql.Open("postgres", testdb.managementURL)
	if err != nil {
		return err
	}

	defer func() {
		if err := sqlDB.Close(); err != nil {
			return
		}
	}()

	_, err = sqlDB.Exec(fmt.Sprintf("DROP DATABASE \"%s\"", testdb.databaseName))

	return err
}
