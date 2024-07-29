package db

import (
	"github.com/gocql/gocql"
	"github.com/zer0day88/cme/message-service/config"
)

func InitCassandra() (*gocql.Session, error) {

	var (
		host     = config.Key.Database.Cassandra.Host
		username = config.Key.Database.Cassandra.Username
		password = config.Key.Database.Cassandra.Password
		port     = config.Key.Database.Cassandra.Port
		keyspace = config.Key.Database.Cassandra.Keyspace
	)

	cluster := gocql.NewCluster(host)
	cluster.Keyspace = keyspace
	cluster.Port = port
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: username,
		Password: password,
	}
	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}

	return session, nil

}
