package main

import (
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
	"net"
	"sync"
	"time"
)

type MockConsulService struct {
	ServiceID   string
	ServiceName string
	Address     string
	ServicePort int
}

// Next steps
// - allow configuring path to Consul endpoint
//		You can do this via: CONSUL_HTTP_ADDR=foo.example.com:8500
// - run repeatedly in a timer
// - later: only print the error output
// - collect up results, sort them by service, then print
// - use a logger

func main() {
	// TODO: Run this repeatedly in a timer
	tryConnectingToConsulServices()
}

func tryConnectingToConsulServices() {
	var wg sync.WaitGroup

	// Lookup services running in Consul
	//services := lookupServices()
	services, err := lookupConsulServices()
	if err != nil {
		fmt.Println("Failed to lookup Consul services. Error:", err)
	}

	// Try a connection to each service instance
	for _, s := range services {
		wg.Add(1)
		go func(s consulapi.CatalogService) {
			// Attempt a connection
			defer wg.Done()

			connString := fmt.Sprintf("%s:%d", s.Address, s.ServicePort)
			fmt.Println("Attempting connection to", connString)
			conn, err := net.DialTimeout("tcp", connString, 5*time.Second)
			if err != nil {
				// TODO: handle error
				fmt.Printf("Failed to connect to: %-20s %-20s %v\n", connString, s.ServiceName, s.ServiceTags)
			} else {
				defer conn.Close()
			}
		}(s)
	}
	wg.Wait()
}

func lookupServices() []consulapi.CatalogService {
	// Some mock services for testing
	tehGoog := consulapi.CatalogService{ServiceID: "broken-google", ServiceName: "google", Address: "google.com", ServicePort: 81}
	google := consulapi.CatalogService{ServiceID: "working-google", ServiceName: "google", Address: "google.com", ServicePort: 80}
	return []consulapi.CatalogService{tehGoog, google}
}

func lookupConsulServices() ([]consulapi.CatalogService, error) {
	// Get a consul client
	client, err := consulapi.NewClient(consulapi.DefaultConfig())
	if err != nil {
		return []consulapi.CatalogService{}, err
	}

	// Lookup catalog of all service names
	catalog := client.Catalog()
	services, _, err := catalog.Services(nil)
	if err != nil {
		return []consulapi.CatalogService{}, err
	}

	// Find all instances of each service
	output := []consulapi.CatalogService{}
	fmt.Println("consul services:")
	for serviceName, _ := range services {
		catalogServices, _, err := catalog.Service(serviceName, "", nil)
		if err != nil {
			return []consulapi.CatalogService{}, err
		}

		fmt.Println(serviceName)
		for _, cs := range catalogServices {
			output = append(output, *cs)
			fmt.Printf("%#v\n", cs)
		}
	}
	return output, nil
}
