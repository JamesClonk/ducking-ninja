package main

type Commands map[string]string

func readCommands(filename string) (commands Commands) {
	readJson(filename, &commands)
	return
}
