package nodes

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type NodeRepository interface {
	Save(node *Node) (err error)
	FindAll() (nodes []*Node, err error)
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
		_ = session.Close()
	}()

	if _, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		return n.persistNode(tx, node)
	}); err != nil {
		return err
	}

	return nil
}

func (n *NodeNeo4jRepository) FindAll() (nodes []*Node, err error) {
	session := n.Driver.NewSession(neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeRead,
	})

	defer func() {
		_ = session.Close()
	}()

	result, err := session.
		ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
			res, err := tx.Run("MATCH (n:Node) RETURN n", nil)

			if err != nil {
				return nil, err
			}

			var nodes []*Node
			for res.Next() {
				record := res.Record()
				if value, ok := record.Get("n"); ok {
					neo4jnode := value.(neo4j.Node)
					props := neo4jnode.Props

					node := &Node{
						Name:            props["name"].(string),
						Alias:           props["alias"].(string),
						IsGateway:       props["gateway"].(bool),
						Platform:        props["platform"].(string),
						OperatingSystem: props["os"].(string),
						Location:        props["location"].(string),
					}
					nodes = append(nodes, node)
				}

			}
			return nodes, nil
		})

	return result.([]*Node), nil
}

func (n *NodeNeo4jRepository) FindByName(name string) (node *Node, err error) {
	session := n.Driver.NewSession(neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeRead,
	})

	defer func() {
		_ = session.Close()
	}()

	result, err := session.
		ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
			res, err := tx.Run("MATCH (n:Node {name: $name}) "+
				"RETURN n.name AS Name, n.alias AS Alias, n.gateway AS IsGateway, "+
				"n.platform AS Platform, n.os AS OperatingSystem, n.location AS Location ",
				map[string]interface{}{
					"name": name,
				})

			if err != nil {
				return nil, err
			}

			singleRecord, err := res.Single()
			if err != nil {
				return nil, err
			}

			return &Node{
				Name:            singleRecord.Values[0].(string),
				Alias:           singleRecord.Values[1].(string),
				IsGateway:       singleRecord.Values[2].(bool),
				Platform:        singleRecord.Values[3].(string),
				OperatingSystem: singleRecord.Values[4].(string),
				Location:        singleRecord.Values[5].(string),
			}, nil
		})

	if result != nil {
		node = result.(*Node)
	} else {
		node = nil
	}

	return node, err
}

func (n *NodeNeo4jRepository) DeleteByName(name string) (err error) {
	session := n.Driver.NewSession(neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})

	defer func() {
		_ = session.Close()
	}()

	_, err = session.
		WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
			res, err := tx.Run("MATCH (n:Node {name: $name}) DELETE n ",
				map[string]interface{}{
					"name": name,
				})

			fmt.Printf("%x", res)
			return nil, err
		})

	return err
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
