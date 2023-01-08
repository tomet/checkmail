return {
	name = "checkmail",
	source_files = {"/*.go"},
	targets = {
		compile = { 
			cmd = "go build .",
			shell = true,
		},
		run = {
			cmd = "go build . && ./checkmail list seen",
			shell = true,
		},
	}
}
 