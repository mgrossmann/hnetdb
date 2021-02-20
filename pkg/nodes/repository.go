package nodes

import "github.com/neo4j/neo4j-go-driver/v4/neo4j"

type NodeRepository interface {
	Save(node *Node) (err error)
	FindAll() (nodes []Node, err error)
	FindByName(name string) (node *Node, err error)
	DeleteByName(name string) (err error)
}

type NodeNeo4jRepository struct {
	Driver neo4j.Driver
}

func (n *NodeNeo4jRepository) Save(node *Node) (err error) {
	session := n.Driver.NewSession(neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})

	defer func() {
		err = session.Close()
	}()

	if _, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		return n.persistNode(tx, node)
	}); err != nil {
		return err
	}

	return nil
}

func (n *NodeNeo4jRepository) FindAll() (nodes []Node, err error) {

	return nil, nil
}

func (n *NodeNeo4jRepository) FindByName(name string) (node *Node, err error) {

	return nil, nil
}

func (n *NodeNeo4jRepository) DeleteByName(name string) (err error) {
	return nil
}

func (n *NodeNeo4jRepository) persistNode(tx neo4j.Transaction, node *Node) (interface{}, error) {

	query := "CREATE (:Node { name: $name, alias: $alias,gateway: $gateway, " +
		"platform: $platform, os: $os, location: $location})"

	parameters := map[string]interface{}{
		"name":     node.Name,
		"alias":    node.Alias,
		"gateway":  node.IsGateway,
		"platform": node.Platform,
		"os":       node.OperatingSystem,
		"location": node.Location,
	}

	_, err := tx.Run(query, parameters)

	return nil, err
}
