package challenge

import _ "embed"

//go:embed challenges.json
var challengesData []byte

//go:embed concepts.json
var conceptsData []byte