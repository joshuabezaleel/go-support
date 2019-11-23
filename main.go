package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strings"

	"github.com/gernest/wow"
	"github.com/gernest/wow/spin"
	"github.com/pkg/browser"
	"github.com/sirkon/goproxy/gomod"
	"github.com/xlab/treeprint"
	"gopkg.in/yaml.v2"
)

const (
	githubAPIURL                  = "https://api.github.com/repos/"
	gitHubSponsorTypeURL          = "https://github.com/sponsors/"
	patreonSponsorTypeURL         = "https://patreon.com/"
	openCollectiveSponsorTypeURL  = "https://opencollective.com/"
	kofiSponsorTypeURL            = "https://ko-fi.com/"
	tideliftSponsorTypeURL        = "https://tidelift.com/"
	communityBridgeSponsorTypeURL = "https://funding.communitybridge.org/projects/"
	liberapaySponsorTypeURL       = "https://en.liberapay.com/"
	issueHuntSponsorTypeURL       = "https://issuehunt.io/r/"
	otechieSponsorTypeURL         = "https://otechie.com/"
)

var (
	moduleGitHubRegex, _      = regexp.Compile("github.com")
	moduleWithVersionRegex, _ = regexp.Compile("github.com/(.*)/(.*)/")
	authorFromAuthorRepo, _   = regexp.Compile("(.*)/")
)

// Sponsor struct reflects sponsorship type supported by GitHub's FUNDING.yml file
type Sponsor struct {
	GitHub          interface{} `yaml:"github"`
	Patreon         string      `yaml:"patreon"`
	OpenCollective  string      `yaml:"open_collective"`
	Kofi            string      `yaml:"ko_fi"`
	Tidelift        string      `yaml:"tidelift"`
	CommunityBridge string      `yaml:"community_bridge"`
	Liberapay       string      `yaml:"liberapay"`
	IssueHunt       string      `yaml:"issuehunt"`
	Otechie         string      `yaml:"otechie"`
	Custom          interface{} `yaml:"custom"`
}

func main() {
	var err error
	authorToProjectsList := make(map[string][]string)
	authorToSponsorsList := make(map[string]Sponsor)
	client := &http.Client{}

	parseResult, err := getModule()
	if err != nil {
		log.Fatal(err)
	}
	if len(parseResult.Require) == 0 {
		fmt.Printf("We didn't find any external Go packages on your project \"%v\" 😶 \n", parseResult.Name)
		os.Exit(1)
	}
	// go loading(parseResult.Name)
	w := wow.New(os.Stdout, spin.Get(spin.Clock), "Retrieving data from GitHub")
	w.Start()

	// Sanitize module not hosted in GitHub and module with version
	var authorRepos []string
	for dependency := range parseResult.Require {
		if moduleGitHubRegex.MatchString(dependency) {
			if moduleWithVersionRegex.MatchString(dependency) {
				authorRepo := moduleWithVersionRegex.FindString(dependency)
				authorRepo = strings.TrimPrefix(authorRepo, "github.com/")
				authorRepo = strings.TrimSuffix(authorRepo, "/")
				authorRepos = append(authorRepos, authorRepo)
			} else {
				authorRepo := strings.TrimPrefix(dependency, "github.com/")
				authorRepos = append(authorRepos, authorRepo)
			}
		}
	}

	for _, authorRepo := range authorRepos {
		req, err := http.NewRequest("GET", githubAPIURL+authorRepo+"/contents/.github/FUNDING.yml", nil)
		req.Header.Set("Authorization", "bearer "+os.Getenv("GITHUB_TOKEN"))
		if err != nil {
			log.Fatal(err)
		}

		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		var respJSON map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&respJSON)
		if err != nil {
			log.Fatal(err)
		}

		if respJSON["content"] != nil {
			decodedContent := fmt.Sprintf("%v", respJSON["content"])
			content, err := base64.StdEncoding.DecodeString(decodedContent)
			if err != nil {
				log.Fatal(err)
			}

			sponsor := Sponsor{}
			err = yaml.Unmarshal(content, &sponsor)
			if err != nil {
				log.Fatal(err)
			}

			// Get author
			var sponsorGitHub string
			switch sponsorGitHubType := sponsor.GitHub.(type) {
			case []interface{}:
				sponsorGitHub = sponsorGitHubType[0].(string)
			case string:
				sponsorGitHub = sponsor.GitHub.(string)
			case nil:
				author := authorFromAuthorRepo.FindString(authorRepo)
				author = strings.TrimSuffix(author, "/")
				sponsorGitHub = author
			}

			_, authorExisted := authorToProjectsList[sponsorGitHub]
			if authorExisted {
				authorToProjectsList[sponsorGitHub] = append(authorToProjectsList[sponsorGitHub], "github.com/"+authorRepo)
			} else {
				authorToSponsorsList[sponsorGitHub] = sponsor
				authorToProjectsList[sponsorGitHub] = []string{"github.com/" + authorRepo}
			}
		}
	}
	w.Stop()

	if len(authorToProjectsList) == 0 {
		fmt.Println("")
		fmt.Printf("We couldn't find any sponsorable Go packages on your project \"%v\" 😞 \n ", parseResult.Name)
		os.Exit(1)
	}

	tree, gitHubURLs := buildTree(authorToProjectsList, authorToSponsorsList)
	fmt.Println("")
	fmt.Println(parseResult.Name)
	fmt.Println(tree.String())

	if len(gitHubURLs) == 0 {
		fmt.Println("We couldn't find any GitHub sponsor pages in the lists but there are still numerous sponsor method available! 🥳")
		os.Exit(1)
	}

	fmt.Print("Do you want to open the GitHub sponsor pages in your browser? (Y/N): ")
	reader := bufio.NewReader(os.Stdin)
	char, _, err := reader.ReadRune()

	if err != nil {
		log.Fatal(err)
	}

	switch char {
	case 'Y', 'y':
		for _, gitHubURL := range gitHubURLs {
			// openbrowser(gitHubURL)
			err = browser.OpenURL(gitHubURL)
			if err != nil {
				log.Fatal(err)
			}
		}
		fmt.Println("Thank you for making the community much better by supporting these awesome Go packages!! ❤️ ❤️ ❤️")
		break
	case 'N':
		fmt.Println("We are looking forward to your support for these awesome Go packages!! 😊")
		break
	}

}

func buildTree(authorToProjectsList map[string][]string, authorToSponsorsList map[string]Sponsor) (tree treeprint.Tree, urls []string) {
	tree = treeprint.New()
	var gitHubURLs []string

	for author, dependencies := range authorToProjectsList {
		authorBranch := tree.AddBranch(author)

		packagesBranch := authorBranch.AddBranch("package(s)")
		for _, dependency := range dependencies {
			packagesBranch.AddNode(dependency)
		}

		sponsorsBranch := authorBranch.AddBranch("donation urls")
		for authorFromList, sponsorType := range authorToSponsorsList {
			if author == authorFromList {
				sponsorStructReflect := reflect.ValueOf(sponsorType)
				for i := 0; i < sponsorStructReflect.NumField(); i++ {
					if sponsorStructReflect.Field(i).Interface() != "" && sponsorStructReflect.Field(i).Interface() != nil {

						sponsorTypeFromReflect := reflect.ValueOf(&sponsorType).Elem().Type().Field(i).Name

						if sponsorTypeFromReflect == "GitHub" {
							switch sponsorGitHubType := sponsorStructReflect.Field(i).Interface().(type) {
							case []interface{}:
								gitHubURLs = append(gitHubURLs, "https://github.com/sponsors/"+sponsorGitHubType[0].(string))
							case string:
								gitHubURLs = append(gitHubURLs, "https://github.com/sponsors/"+sponsorGitHubType)
							}
						}

						switch sponsorReflectType := sponsorStructReflect.Field(i).Interface().(type) {
						case []interface{}:
							if len(sponsorReflectType) == 1 {
								sponsorsBranch.AddNode(sponsorTypeFromReflect + ": " + appendSponsorTypeURL(sponsorTypeFromReflect, sponsorReflectType[0].(string)))
							} else {
								sponsorTypeBranch := sponsorsBranch.AddBranch(sponsorTypeFromReflect)
								for _, sponsorURL := range sponsorReflectType {
									sponsorTypeBranch.AddNode(appendSponsorTypeURL(sponsorTypeFromReflect, sponsorURL.(string)))
								}
							}
						case string:
							sponsorsBranch.AddNode(sponsorTypeFromReflect + ": " + appendSponsorTypeURL(sponsorTypeFromReflect, sponsorStructReflect.Field(i).Interface().(string)))
						}
					}
				}
			}
		}
	}

	return tree, gitHubURLs
}

func getModule() (*gomod.Module, error) {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	input, err := ioutil.ReadFile(dir + "/go.mod")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	fileName := ""
	parseResult, err := gomod.Parse(fileName, input)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return parseResult, nil
}

func appendSponsorTypeURL(sponsorType, sponsor string) string {
	var appendedSponsor string
	switch sponsorType {
	case "GitHub":
		appendedSponsor = gitHubSponsorTypeURL + sponsor
	case "Patreon":
		appendedSponsor = patreonSponsorTypeURL + sponsor
	case "OpenCollective":
		appendedSponsor = openCollectiveSponsorTypeURL + sponsor
	case "Kofi":
		appendedSponsor = kofiSponsorTypeURL + sponsor
	case "Tidelift":
		appendedSponsor = tideliftSponsorTypeURL + sponsor
	case "ComunityBridge":
		appendedSponsor = communityBridgeSponsorTypeURL + sponsor
	case "Liberapay":
		appendedSponsor = liberapaySponsorTypeURL + sponsor
	case "IssueHunt":
		appendedSponsor = issueHuntSponsorTypeURL + sponsor
	case "Otechie":
		appendedSponsor = otechieSponsorTypeURL + sponsor
	case "Custom":
		appendedSponsor = sponsor
	}

	return appendedSponsor
}

// func openbrowser(url string) {
// 	var err error

// 	switch runtime.GOOS {
// 	case "linux":
// 		err = exec.Command("xdg-open", url).Start()
// 		break
// 	case "windows":
// 		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
// 		break
// 	case "darwin":
// 		err = exec.Command("open", url).Start()
// 		break
// 	default:
// 		err = fmt.Errorf("unsupported platform")
// 	}
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// }
