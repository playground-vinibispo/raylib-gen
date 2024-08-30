/*
Copyright © 2024 Vinícius Bispo vini.bispo015@gmail.com
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "raylib-gen",
	Short: "Generate code for raylib in Go",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func askForProjectName() string {
	var name string
	fmt.Print("What is the name of the project? ")
	fmt.Scanln(&name)
	return name
}

func createFolder(name string) {
	os.Mkdir(name, 0755)
	os.Chdir(name)
}

func runGoMod(name string) {
	fmt.Println("Preparing go.mod file...")
	exec.Command("go", "mod", "init", name).Run()
}

func installRaylib() {
	fmt.Println("Installing raylib-go...")
	exec.Command("go", "get", "-v", "-u", "github.com/gen2brain/raylib-go/raylib").Run()
}

func initGit() {

	if _, err := exec.LookPath("git"); err == nil {
		fmt.Println("Initializing git repository...")
		exec.Command("git", "init").Run()
		fmt.Println("Checking out to main branch...")
		exec.Command("git", "branch", "-M", "main").Run()
	}
}

func createGitIgnore() {
	if _, err := exec.LookPath("npx"); err == nil {
		fmt.Println("Creating .gitignore file...")
		exec.Command("npx", "gitignore", "go").Run()
	}
}

func runAirInit() {
	if _, err := exec.LookPath("air"); err == nil {
		fmt.Println("Initializing air...")
		exec.Command("air", "init").Run()
	}
}

func createGoFiles(name string) {

	os.Mkdir("internals", 0755)
	fmt.Println("Creating cmd folder...")
	os.MkdirAll(fmt.Sprintf("cmd/%s", name), 0755)
	fmt.Println("Creating main.go file...")
	file, err := os.Create(fmt.Sprintf("cmd/%s/main.go", name))
	if err != nil {
		fmt.Println("Error creating main.go file")
		os.Exit(1)
	}
	fmt.Println("Writing to main.go file...")
	_, err = file.WriteString(fmt.Sprintf(`package main
    import (
      rl "github.com/gen2brain/raylib-go/raylib"
    )
    func main() {
      rl.InitWindow(400, 400, "%s")
      rl.SetTargetFPS(60)
      for !rl.WindowShouldClose() {
        rl.BeginDrawing()
        rl.ClearBackground(rl.RayWhite)
        rl.DrawText("Hello, world!", 12, 12, 20, rl.Maroon)
        rl.EndDrawing()
      }
      rl.CloseWindow()
    }
    `, name))
	if err != nil {
		fmt.Println("Error writing to main.go file")
		os.Exit(1)
	}
	defer file.Close()
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new project",
	Run: func(cmd *cobra.Command, args []string) {
		name := askForProjectName()
		createFolder(name)
		runGoMod(name)
		installRaylib()
		initGit()
		createGitIgnore()
		createGoFiles(name)
		runAirInit()

		if _, err := exec.LookPath("air"); err != nil {
			fmt.Println("Air is not installed. Please install it to continue")
			os.Exit(1)
		}
		fmt.Println("Editing .air.toml file...")
		byteContent, err := os.ReadFile(".air.toml")
		if err != nil {
			fmt.Println(err)
			fmt.Println("Error reading .air.toml file")
			os.Exit(1)
		}
		content := string(byteContent)
		re := regexp.MustCompile(`cmd = ".*"`)
		newCmdValue := "go build -o ./tmp/main cmd/**/main.go"
		newContent := re.ReplaceAllString(content, fmt.Sprintf(`cmd = "%s"`, newCmdValue))
		fmt.Println("New content")
		file, err := os.OpenFile(".air.toml", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			fmt.Println("Error opening .air.toml file")
			os.Exit(1)
		}
		defer file.Close()
		if _, err = file.WriteString(newContent); err != nil {
			fmt.Println(err)
			fmt.Println("Error writing to .air.toml file")
			os.Exit(1)
		}

		fmt.Println("Project initialized")
		exec.Command("go", "mod", "tidy").Run()
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
