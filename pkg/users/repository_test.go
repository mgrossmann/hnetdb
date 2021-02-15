package users

import (
	"context"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/testcontainers/testcontainers-go"
)

var _ = Describe("User repository", func() {

	It("register users", func() {
		ctx := context.Background()

		request := testcontainers.ContainerRequest{
			Image:        "neo4j",
			ExposedPorts: []string{"7687/tcp"},
		}

		neo4j, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
			ContainerRequest: request,
			Started:          true,
		})

		Expect(err).To(BeNil(), "NEO4J container should start")
		defer neo4j.Terminate(ctx)

		host, err := neo4j.Host(ctx)
		Expect(err).To(BeNil(), "Should get th host name")

		port, err := neo4j.MappedPort(ctx, "7687")
		Expect(err).To(BeNil(), "Should get the port")

		fmt.Printf("NEO4J was started ad %s:%d", host, port.Int())
	})

})