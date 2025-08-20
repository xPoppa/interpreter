testf: 
	air -build.cmd "go test -v ./... | grep -i fail -C 3"
