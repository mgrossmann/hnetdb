package main

import (
	"github.com/mvslovers/hnetdb/pkg/nodes"
	"github.com/mvslovers/hnetdb/pkg/users"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"net/http"
	"os"
)

func main() {

	neo4jUri, found := os.LookupEnv("NEO4J_URI")
	if !found {
		panic("NEO4J_URI not set")
	}
	neo4jUsername, found := os.LookupEnv("NEO4J_USERNAME")
	if !found {
		panic("NEO4J_USERNAME not set")
	}
	neo4jPassword, found := os.LookupEnv("NEO4J_PASSWORD")
	if !found {
		panic("NEO4J_PASSWORD not set")
	}
	usersRepository := users.UserNeo4jRepository{
		Driver: driver(neo4jUri, neo4j.BasicAuth(neo4jUsername, neo4jPassword, "")),
	}
	nodesRepository := nodes.NodeNeo4jRepository{
		Driver: driver(neo4jUri, neo4j.BasicAuth(neo4jUsername, neo4jPassword, "")),
	}

	registrationHandler := &users.UserRegistrationHandler{
		Path:           "/users/register",
		UserRepository: &usersRepository,
	}
	loginHandler := &users.UserLoginHandler{
		Path:           "/users/login",
		UserRepository: &usersRepository,
	}
	newNodeHandler := &nodes.NewNodeHandler{
		Path:           "/node",
		NodeRepository: &nodesRepository,
	}
	server := http.NewServeMux()
	server.HandleFunc(registrationHandler.Path, registrationHandler.Register)
	server.HandleFunc(loginHandler.Path, loginHandler.Login)
	server.HandleFunc(newNodeHandler.Path, newNodeHandler.New)

	if err := http.ListenAndServe(":3000", server); err != nil {
		panic(err)
	}
}

func driver(target string, token neo4j.AuthToken) neo4j.Driver {
	result, err := neo4j.NewDriver(target, token)
	if err != nil {
		panic(err)
	}
	return result
}
