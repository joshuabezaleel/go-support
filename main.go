package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/sirkon/goproxy/gomod"
)

func main() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	input, err := ioutil.ReadFile(dir + "/go.mod")
	if err != nil {
		log.Fatal(err)
	}

	fileName := ""
	parseResult, err := gomod.Parse(fileName, input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(parseResult.Name)

	authorToDependencies := make(map[string][]string)
	for dependency := range parseResult.Require {
		for author, projects := range authorToProjectsList {
			for _, project := range projects {
				if dependency == project {
					found := false
					for existingAuthor := range authorToDependencies {
						if existingAuthor == author {
							found = true
							// tambah ke []string dengan key existingAuthor
							authorToDependencies[author] = append(authorToDependencies[author], dependency)
						}
					}
					if !found {
						var dependencies = []string{dependency}
						authorToDependencies[author] = dependencies
					}
				}
			}
		}
	}

	for author, dependencies := range authorToDependencies {
		fmt.Println(author)
		for _, dependency := range dependencies {
			fmt.Println(dependency)
		}
	}

}
