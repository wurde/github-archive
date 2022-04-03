package main

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/joho/godotenv"
)

var TIME_NOW = time.Now().Format("20060102")

func deleteEmpty(arr []string) []string {
	var newArr []string
	for _, str := range arr {
		if str != "" {
			newArr = append(newArr, str)
		}
	}
	return newArr
}

func main() {
	log.Println("Running GitHub Archive...")

	// Load local environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Check for a GitHub personal access token
	var GITHUB_TOKEN = os.Getenv("GITHUB_TOKEN")
	if GITHUB_TOKEN == "" {
		log.Fatal("GITHUB_TOKEN is not set")
	}

	// Read list of repositories to archive
	content, err := os.ReadFile("./repos.txt")
	if err != nil {
		log.Fatal(err)
	}
	repos := strings.Split(string(content), "\n")
	repos = deleteEmpty(repos)
	if len(repos) <= 0 {
		log.Fatal("No repositories found in ./repos.txt")
	}

	// If repos already exist, remove them
	os.RemoveAll("./tmp")

	// Clone repos locally to ./tmp
	for _, repo := range repos {
		log.Printf("Cloning %s...", repo)

		re := regexp.MustCompile(`([\w-]+)$`)
		repoName := re.FindString(repo)
		repoPath := "./tmp/" + repoName

		_, err = git.PlainClone(repoPath, false, &git.CloneOptions{
			URL: "https://" + GITHUB_TOKEN + "@" + repo,
		})
		if err != nil {
			log.Fatal(err)
		}

		// Zip the repo ./tmp/<repo>-YYYYMMDD.zip
		outFile, err := os.Create(fmt.Sprintf("./tmp/%s-%s.zip", repoName, TIME_NOW))
		if err != nil {
			log.Fatal(err)
		}
		defer outFile.Close()
		myZip := zip.NewWriter(outFile)
		defer myZip.Close()

		err = filepath.Walk(repoPath, func(filePath string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			if err != nil {
				return err
			}
			relPath := strings.TrimPrefix(filePath, filepath.Dir(repoPath))
			zipFile, err := myZip.Create(relPath)
			if err != nil {
				return err
			}
			fsFile, err := os.Open(filePath)
			if err != nil {
				return err
			}
			_, err = io.Copy(zipFile, fsFile)
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			log.Fatal(err)
		}
		err = myZip.Close()
		if err != nil {
			log.Fatal(err)
		}
	}
}
