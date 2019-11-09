package main

var authorToProjectsList = map[string][]string{
	"tj": []string{
		"github.com/tj/go-terminput",
		"github.com/tj/go-editor",
	},
	"deadprogram": []string{
		"github.com/hybridgroup/gocv",
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
}

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
}
