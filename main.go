package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"

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

	tree, urls := buildTree(authorToDependencies)
	fmt.Println(tree.String())

	fmt.Print("Do you want to open the donation pages in browser? (Y/N): ")
	reader := bufio.NewReader(os.Stdin)
	char, _, err := reader.ReadRune()

	if err != nil {
		log.Fatal(err)
	}

	switch char {
	case 'Y':
		for _, url := range urls {
			openbrowser(url)
		}
		fmt.Println("Thank you for supporting these awesome Go packages!! :)")
		break
	case 'N':
		fmt.Println("We are looking forward for your support for these awesome Go packages!! :)")
		break
	}

}

func buildTree(authorToDependencies map[string][]string) (tree treeprint.Tree, urls []string) {
	tree = treeprint.New()

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
					urls = append(urls, url)
				}
			}
		}
	}

	return tree, urls
}

func openbrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
		break
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
		break
	case "darwin":
		err = exec.Command("open", url).Start()
		break
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}

}
