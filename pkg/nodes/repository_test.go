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
