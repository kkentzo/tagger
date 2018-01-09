package main

func main() {
	project := NewProject(
		"/Users/kkentzo/Workspace/agnostic_backend",
		"",
		[]string{".git", "coverage"})

	project.Monitor()

}
