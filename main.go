package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/xlab/treeprint"

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

	tree := buildTree(authorToDependencies)
	fmt.Println(tree.String())

}

func buildTree(authorToDependencies map[string][]string) treeprint.Tree {
	tree := treeprint.New()

	for author, dependencies := range authorToDependencies {
		authorBranch := tree.AddBranch(author)
		packagesBranch := authorBranch.AddBranch("package(s)")
		for _, dependency := range dependencies {
			packagesBranch.AddNode(dependency)
		}
		sponsorBranch := authorBranch.AddBranch("donation urls")
		for authorFromList, urlsFromList := range authorToSponsorURLsList {
			if author == authorFromList {
				for urlType, url := range urlsFromList {
					donationURL := urlType + ": " + url
					sponsorBranch.AddNode(donationURL)
				}
			}
		}
	}

	return tree
}
