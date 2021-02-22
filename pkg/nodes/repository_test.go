package nodes_test

import (
	"context"
	"fmt"
	. "github.com/mvslovers/hnetdb/pkg/nodes"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"io"
)

var _ = Describe("Testing node repository", func() {

	const username = "neo4j"
	const password = "s3cr3t"

	var ctx context.Context

	var neo4jContainer testcontainers.Container
	var driver neo4j.Driver

	var repository NodeRepository

	BeforeSuite(func() {
		ctx = context.Background()
		var err error

		neo4jContainer, err = startContainer(ctx, username, password)
		Expect(err).To(BeNil(), "Container should start")

		port, err := neo4jContainer.MappedPort(ctx, "7687")
		Expect(err).To(BeNil(), "Port should be resolved")

		address := fmt.Sprintf("bolt://localhost:%d", port.Int())

		driver, err = neo4j.NewDriver(address, neo4j.BasicAuth(username, password, ""))
		Expect(err).To(BeNil(), "Driver should be created")

		repository = &NodeNeo4jRepository{
			Driver: driver,
		}
	})

	AfterSuite(func() {
		Close(driver, "Driver")
		Expect(neo4jContainer.Terminate(ctx)).To(BeNil(), "Container should stop")
	})

	It("Save", func() {
		newNode := &Node{
			Name:            "DRNBRX1A",
			Alias:           "DEBRXMVS",
			IsGateway:       false,
			Platform:        "Hercules 4 on Linux",
			OperatingSystem: "MVS3.8J",
			Location:        "Germany",
		}

		// test repository
		err := repository.Save(newNode)

		Expect(err).To(BeNil(), "Node should be created")

		// read the created node for test
		session := driver.NewSession(neo4j.SessionConfig{})

		defer Close(session, "Session")

		result, err := session.
			ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
				res, err := tx.Run("MATCH (n:Node {name: $name}) "+
					"RETURN n.name AS Name, n.alias AS Alias, n.gateway AS IsGateway, "+
					"n.platform AS Platform, n.os AS OperatingSystem, n.location AS Location ",
					map[string]interface{}{
						"name": newNode.Name,
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

		Expect(err).To(BeNil(), "The read transaction should not fail")

		testNode := result.(*Node)

		Expect(testNode.Name).To(Equal(newNode.Name))
		Expect(testNode.Alias).To(Equal(newNode.Alias))
		Expect(testNode.IsGateway).To(Equal(newNode.IsGateway))
		Expect(testNode.Platform).To(Equal(newNode.Platform))
		Expect(testNode.OperatingSystem).To(Equal(newNode.OperatingSystem))
		Expect(testNode.Location).To(Equal(newNode.Location))

	})

	It("FindAll", func() {
		testNode1 := &Node{
			Name:            "DRNMIG1A",
			Alias:           "DEMIGVM",
			IsGateway:       false,
			Platform:        "Hercules 4 on Linux",
			OperatingSystem: "VM/ESA",
			Location:        "Germany",
		}

		testNode2 := &Node{
			Name:            "DRNMIG3A",
			Alias:           "DEMIGMVS",
			IsGateway:       false,
			Platform:        "Hercules 4 on Linux",
			OperatingSystem: "MVS3.8J",
			Location:        "Germany",
		}

		// create two nodes for testing
		session := driver.NewSession(neo4j.SessionConfig{})

		defer Close(session, "Session")

		// test node #1
		_, err := session.
			WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
				query := "CREATE (:Node { name: $name, alias: $alias,gateway: $gateway, " +
					"platform: $platform, os: $os, location: $location})"

				parameters := map[string]interface{}{
					"name":     testNode1.Name,
					"alias":    testNode1.Alias,
					"gateway":  testNode1.IsGateway,
					"platform": testNode1.Platform,
					"os":       testNode1.OperatingSystem,
					"location": testNode1.Location,
				}

				_, err := tx.Run(query, parameters)

				return nil, err
			})

		Expect(err).To(BeNil(), "The write transaction should not fail")

		// test node #2
		_, err = session.
			WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
				query := "CREATE (:Node { name: $name, alias: $alias,gateway: $gateway, " +
					"platform: $platform, os: $os, location: $location})"

				parameters := map[string]interface{}{
					"name":     testNode2.Name,
					"alias":    testNode2.Alias,
					"gateway":  testNode2.IsGateway,
					"platform": testNode2.Platform,
					"os":       testNode2.OperatingSystem,
					"location": testNode2.Location,
				}

				_, err := tx.Run(query, parameters)

				return nil, err
			})

		Expect(err).To(BeNil(), "The write transaction should not fail")

		// test repository function FindByAll
		allNodes, err := repository.FindAll()

		var foundNode1 *Node
		var foundNode2 *Node

		for _, node := range allNodes {
			Expect(node).To(Not(BeNil()), "First node found, should not be nil")

			if node.Name == testNode1.Name {
				foundNode1 = node
			}

			if node.Name == testNode2.Name {
				foundNode2 = node
			}
		}

		Expect(foundNode1).To(Equal(testNode1), "We should have found our 1st testing node")
		Expect(foundNode2).To(Equal(testNode2), "We should have found our 2nd testing node")
	})

	It("FindByName", func() {
		testNode := &Node{
			Name:            "DRNBRX9A",
			Alias:           "DEBRXLNX",
			IsGateway:       false,
			Platform:        "Linux",
			OperatingSystem: "Linux",
			Location:        "Germany",
		}

		// create a node for testing
		session := driver.NewSession(neo4j.SessionConfig{})
		defer Close(session, "Session")

		_, err := session.
			WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
				query := "CREATE (:Node { name: $name, alias: $alias,gateway: $gateway, " +
					"platform: $platform, os: $os, location: $location})"

				parameters := map[string]interface{}{
					"name":     testNode.Name,
					"alias":    testNode.Alias,
					"gateway":  testNode.IsGateway,
					"platform": testNode.Platform,
					"os":       testNode.OperatingSystem,
					"location": testNode.Location,
				}

				_, err := tx.Run(query, parameters)

				return nil, err
			})
		Expect(err).To(BeNil(), "The write transaction should not fail")

		// test repository function FindByName
		foundNode, err := repository.FindByName(testNode.Name)
		Expect(err).To(BeNil(), "FindByeName should not end with an error")
		Expect(foundNode).To(Not(BeNil()), "Node should be found")
		Expect(foundNode).To(Equal(testNode))

		// test repository function FindByName with a non existing node
		foundNode, err = repository.FindByName("DUMMY")
		Expect(err).To(Not(BeNil()), "Should end with an error")
	})

	It("DeleteByName", func() {
		testNode := &Node{
			Name:            "DRNMIG9A",
			Alias:           "DEMIGLNX",
			IsGateway:       false,
			Platform:        "Linux",
			OperatingSystem: "Linux",
			Location:        "Germany",
		}

		// create a node for testing
		session := driver.NewSession(neo4j.SessionConfig{})
		defer Close(session, "Session")

		_, err := session.
			WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
				query := "CREATE (:Node { name: $name, alias: $alias,gateway: $gateway, " +
					"platform: $platform, os: $os, location: $location})"

				parameters := map[string]interface{}{
					"name":     testNode.Name,
					"alias":    testNode.Alias,
					"gateway":  testNode.IsGateway,
					"platform": testNode.Platform,
					"os":       testNode.OperatingSystem,
					"location": testNode.Location,
				}

				_, err := tx.Run(query, parameters)

				return nil, err
			})
		Expect(err).To(BeNil(), "The write transaction should not fail")

		// test repository function DeleteByName
		err = repository.DeleteByName(testNode.Name)
		Expect(err).To(BeNil(), "Node should be deleted")

		// test if node is really gone
		_, err = session.
			ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
				res, err := tx.Run("MATCH (n:Node {name: $name}) "+
					"RETURN n.name AS Name, n.alias AS Alias, n.gateway AS IsGateway, "+
					"n.platform AS Platform, n.os AS OperatingSystem, n.location AS Location ",
					map[string]interface{}{
						"name": testNode.Name,
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
		Expect(err).To(Not(BeNil()), "Should end with an error")

	})

})

func Close(closer io.Closer, resourceName string) {
	Expect(closer.Close()).
		To(BeNil(), "%s should close", resourceName)
}

func startContainer(ctx context.Context, username, password string) (testcontainers.Container, error) {
	request := testcontainers.ContainerRequest{
		Image:        "neo4j",
		ExposedPorts: []string{"7687/tcp"},
		Env:          map[string]string{"NEO4J_AUTH": fmt.Sprintf("%s/%s", username, password)},
		WaitingFor:   wait.ForLog("Bolt enabled"),
	}
	return testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: request,
		Started:          true,
	})
}
