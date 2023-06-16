package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"strings"

	"github.com/joho/godotenv"
)

const dtFormat string = "2006-01-02 15:04:05 Monday"
const logsBreakLine string = "------------------------ mongodb-backup-script ------------------------"

func Handle(err error) {
	if err != nil {
		log.Panic(err.Error())
	}
}

// Basic logging
func LogToFile(message string, debug bool) {
	path, _ := os.Getwd()
	var logFile string = "databaseBackup.log"
	if debug {
		logFile = "databaseBackup.debug.log"
	}
	f, err := os.OpenFile(filepath.Join(path, logFile), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	Handle(err)
	defer f.Close()

	log.SetOutput(f)
	log.Println(fmt.Sprintf(" | %s", message))
}

func UploadToGithub(archiveDir string) {
	// Add archive files to stage
	cmd := exec.Command("git", "add", ".")
	cmd.Dir = archiveDir
	stdout, err := cmd.Output()
	LogToFile(string(stdout), true)
	Handle(err)

	now := time.Now()

	var author string = os.Getenv("github_author")
	var email string = os.Getenv("github_email")

	if author == "" || email == "" {
		LogToFile(fmt.Sprintf("ERROR - github link provided but no author or email, add to .env: github_author=<username> and github_email=<email>\n%s", logsBreakLine), true)
		LogToFile(fmt.Sprintf("ERROR - github link provided but no author or email, add to .env: github_author=<username> and github_email=<email>\n%s", logsBreakLine), false)
		log.Panic("Missing required .env: github_author, github_email")
	}

	// Commit archive files
	cmd = exec.Command("git", "commit", "-m", fmt.Sprintf("'%s'", now.Format(dtFormat)), fmt.Sprintf("--author=\"%s <%s>\"", author, email))
	cmd.Dir = archiveDir
	stdout, err = cmd.Output()
	LogToFile(string(stdout), true)
	Handle(err)

	// Push archive files to repository
	cmd = exec.Command("git", "push", "origin", "master", "--force")
	cmd.Dir = archiveDir
	stdout, err = cmd.Output()
	LogToFile(string(stdout), true)
	Handle(err)

	LogToFile("Successfully uploaded archive to github repository", false)
}

func main() {
	path, _ := os.Getwd()

	// Load environment variables
	godotenv.Load(filepath.Join(path, ".env"))
	var mongoURI string = os.Getenv("mongoURI")
	var database_string string = os.Getenv("databases")
	var github string = os.Getenv("github")

	// Ensure required variables
	if mongoURI == "" || database_string == "" {
		LogToFile(fmt.Sprintf("ERROR - Missing required .env variables: mongoURI, databases\n%s", logsBreakLine), true)
		LogToFile(fmt.Sprintf("ERROR - Missing required .env variables: mongoURI, databases\n%s", logsBreakLine), false)
		log.Panic("Missing required .env variables: mongoURI, databases")
	}

	var databases []string = strings.Split(database_string, ", ")

	// Ensure archive directory exists // initialize repository if github provided
	archiveDir := filepath.Join(path, "archive")
	if _, err := os.Stat(archiveDir); os.IsNotExist(err) {
		err := os.Mkdir(archiveDir, os.ModePerm)
		Handle(err)

		if github != "" {
			cmd := exec.Command("git", "init")
			cmd.Dir = archiveDir
			stdout, err := cmd.Output()
			LogToFile(string(stdout), true)
			Handle(err)

			cmd = exec.Command("git", "remote", "add", "origin", github)
			cmd.Dir = archiveDir
			stdout, err = cmd.Output()
			LogToFile(string(stdout), true)
			Handle(err)
		}
	}

	// Perform mongodump to archive databases to .gzip format
	for _, db := range databases {
		archivePath := filepath.Join(archiveDir, db+".gzip")
		cmd := exec.Command(
			"mongodump",
			"--uri="+mongoURI,
			"--authenticationDatabase=admin",
			"--db="+db,
			"--archive="+archivePath,
			"--gzip",
		)
		stdout, err := cmd.Output()
		LogToFile(string(stdout), true)
		Handle(err)

		LogToFile(fmt.Sprintf("Successfully archived %s", db), false)
	}

	if github != "" {
		UploadToGithub(archiveDir)
	}

	LogToFile(logsBreakLine, false) // Log break for easier viewing
	LogToFile(logsBreakLine, true)  // Log break for easier viewing
}
