package main

// authorToProjectList is a map with the string as key and slice of string as value.
// The key is the name of the author while the slice of string is the list of the projects of an author.
// FOR DEVELOPERS
// List your name and projects conforming to the previous examples with the template of:
// "$your-name": []string{
// 	"github.com/$your-name/$your-awesome-project-1",
// 	"github.com/$your-name/$your-awesome-project-2",
// },
var authorToProjectsList = map[string][]string{
	"tj": []string{
		"github.com/tj/go-terminput",
		"github.com/tj/go-editor",
	},
	"deadprogram": []string{
		"github.com/hybridgroup/gocv",
		"github.com/hybridgroup/gobot",
	},
	"markbates": []string{
		"github.com/gobuffalo/buffalo",
	},
	"mholt": []string{
		"github.com/caddyserver/caddy",
	},
	"colly": []string{
		"github.com/gocolly/colly",
	},
	"peterbourgon": []string{
		"github.com/go-kit/kit",
	},
	"micro": []string{
		"github.com/micro/go-plugins",
		"github.com/micro/micro",
		"github.com/micro/go-micro",
		"github.com/micro/macro",
		"github.com/micro/cli",
	},
}

// authorToSponsorURLsList is a map with the string as key and a map[string]string as value.
// The first string key is the name of the author with the value of map[string]string.
// The second string key is the type of the donation page (github, patreon, opencollective, personal site, etc.)
// while the value is the corresponding URL of the page.
// FOR DEVELOPERS
// List your name and your donation page URSL conforming to the previous examples with the template of:
// "$your-name": map[string]string{
// 	"patreon": "https://patreon.com/$your-name",
// 	"github.com": "https://github.com/sponrs/$your-name",
//  "opencollective": "https://opencollective.com/$your-name",
// },
var authorToSponsorURLsList = map[string]map[string]string{
	"tj": map[string]string{
		"github": "https://github.com/sponsors/tj",
	},
	"deadprogram": map[string]string{
		"patreon": "https://www.patreon.com/deadprogram",
	},
	"markbates": map[string]string{
		"patreon": "https://patreon.com/buffalo",
	},
	"mholt": map[string]string{
		"github": "https://github.com/sponsors/mholt",
	},
	"colly": map[string]string{
		"opencollective": "https://opencollective.com/colly",
	},
	"peterbourgon": map[string]string{
		"github": "https://github.com/sponsors/peterbourgon",
	},
	"micro": map[string]string{
		"issuehunt": "https://issuehunt.io/r/micro/development",
	},
}
