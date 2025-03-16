package main

import (
	"github.com/mihai-valentin/polyroll/internal"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 || os.Args[1] == "" {
		log.Fatalln("missing required argument 1 - path to config yaml file")
	}

	config, err := internal.ReadConfigFromFile(os.Args[1])
	if err != nil {
		log.Fatalf("error reading config file: %s", err)
	}

	ec := internal.NewElkClient(config.ElkHost, config.AuthToken)

	for _, policy := range config.IlmPolicies {
		log.Printf("Creating policy [%s]...\n", policy.Name)

		if err := ec.CreateOrUpdateIlmPolicy(policy); err != nil {
			log.Printf("Cannot create ILM policy [%s]: %s\n", policy.Name, err)
			continue
		}

		log.Printf("Successfully created ILM policy [%s]\n", policy.Name)
	}

	for _, indexTemplate := range config.IndexTemplates {
		log.Printf("Creating index template [%s] with ILM policy [%s]...\n",
			indexTemplate.Name,
			indexTemplate.IlmPolicyName,
		)

		if err := ec.CreateOrUpdateIndexTemplate(indexTemplate); err != nil {
			log.Printf("Cannot create index template [%s]: %s\n", indexTemplate.Name, err)
			continue
		}

		log.Printf("Successfully created index template [%s]\n", indexTemplate.Name)
	}
}
