package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func hashFile(file_path string) (string, error) {
	file_Content, err := os.Open(file_path)
	if err != nil {
		return "", err
	}
	defer file_Content.Close()
	hash := sha256.New()
	if _, err := io.Copy(hash, file_Content); err != nil {
		return "", err
	}
	hashInBytes := hash.Sum(nil)
	return fmt.Sprintf("%x", hashInBytes), nil
}

func findHashDupes(folderDirs []string) (map[string][]string, error) {
	hashesOfFiles := make(map[string][]string)

	for _, folderDir := range folderDirs {
		err := filepath.Walk(folderDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				hash, err := hashFile(path)
				if err != nil {
					return err
				}
				hashesOfFiles[hash] = append(hashesOfFiles[hash], path)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	return hashesOfFiles, nil
}

func deleteDupes(duplicates map[string][]string, keepOriginal bool) error {
	for _, files := range duplicates {
		if len(files) > 1 {

			var latestFile string
			var latestModTime time.Time

			for _, file := range files {
				info, err := os.Stat(file)
				if err != nil {
					fmt.Println("yo what the heck, can't even get file info for", file, err)
					continue
				}

				if info.ModTime().After(latestModTime) {
					latestFile = file
					latestModTime = info.ModTime()
				}
			}

			for _, file := range files {
				if file != latestFile {
					fmt.Println("\noops, another duplicate... deleting: ", file)
					err := os.Remove(file)
					if err != nil {
						fmt.Println("lol couldn't delete that one: ", file, err)
						continue
					}
					fmt.Println("byeee, boomed: ", file)
				} else if keepOriginal {
					fmt.Println("keeping the freshest one, congrats: ", file)
				}
			}
		}
	}
	return nil
}

func deleteAllDupes(duplicates map[string][]string) error {
	for _, files := range duplicates {
		if len(files) > 1 {
			for _, file := range files {
				fmt.Println("\noops, say hello to my little friend: ", file)
				err := os.Remove(file)
				if err != nil {
					fmt.Println("lol couldn't delete that one: ", file, err)
					continue
				}
				fmt.Println("just boomed: ", file)
			}
		}
	}
	return nil
}

func main() {
	fmt.Println("yo, let me cook")

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("uhh... senior dev, can’t find user dir, help pls:", err)
		return
	}

	folderDirs := []string{
		filepath.Join(homeDir, "downloads"),
	}

	duplicates, err := findHashDupes(folderDirs)
	if err != nil {
		fmt.Println("bruh, error finding dupes:", err)
		return
	}

	if len(duplicates) > 0 {
		for hash, files := range duplicates {
			if len(files) > 1 {
				fmt.Println("\nyo, these files share the same mama (hash): ", hash)
				for _, file := range files {
					fmt.Println("here’s the nightmare: ", file)
				}
			}
		}

		fmt.Println("\nso... what’s the move?")
		fmt.Println("1. delete everything, no one’s special (no files kept).")
		fmt.Println("2. delete the dupes and keep the freshest one.")
		fmt.Println("3. skip it, you just wanna see the hashes.")

		var choice string
		fmt.Scanln(&choice)

		switch strings.TrimSpace(choice) {
		case "1":
			err := deleteAllDupes(duplicates)
			if err != nil {
				fmt.Println("yo, couldn’t delete those dupes, wtf:", err)
			}
		case "2":
			err := deleteDupes(duplicates, true)
			if err != nil {
				fmt.Println("yo, couldn’t delete those dupes, wtf:", err)
			}
		case "3":
			fmt.Println("ight bud, skipping deletes. you just wanted the dough.")
		default:
			fmt.Println("yo, that’s not an option, but whatever. skipping anyway.")
		}
	} else {
		fmt.Println("no dupes found, you’re clear, no worries.")
	}
}
