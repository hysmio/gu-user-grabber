package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/pkg/browser"
)

var prefix = []byte("apolloId: ")
var ourID = ""

var ErrNotFound = errors.New("id not found")

type Config struct {
	FilePath string `json:"file_path"`
}

func ReadCfg() Config {
	cfg := Config{}

	inferredDirectory, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln(fmt.Errorf("couldn't infer the log directory: %s", err))
	}

	cfg.FilePath = inferredDirectory + `\AppData\LocalLow\Immutable\gods\debug.log`

	file, err := os.Open("config.json")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			file, err = os.Create("config.json")
			if err != nil {
				log.Fatalln(fmt.Errorf("config.json did not exist, couldn't create it: %s", err))
			}
			d, _ := json.Marshal(cfg)
			_, err = file.WriteString(string(d))
			if err != nil {
				log.Fatalln(fmt.Errorf("couldn't save config to config.json: %s", err))
			}

			file.Seek(0, 0)
		} else {
			log.Fatalln(fmt.Errorf("couldn't open config.json: %s", err))
		}
	}

	err = json.NewDecoder(file).Decode(&cfg)
	if err != nil {
		if errors.Is(err, io.EOF) {
			log.Fatalln(fmt.Errorf("file is empty, if you're unsure of the path, delete the file and it will recreate it"))
		}
		log.Fatalln(fmt.Errorf("couldn't parse config.json as valid json: %s", err))
	}

	return cfg
}

func main() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalln(fmt.Errorf("couldn't create a file watcher: %s", err))
	}

	cfg := ReadCfg()
	fmt.Printf("attempting to watch %s, if this is not the valid log file, i will wait indefinintely, if nothing happens when the game launches, this is the wrong path\n", cfg.FilePath)

	for {
		err = watcher.Add(cfg.FilePath)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				time.Sleep(time.Second)
				continue
			}
			log.Fatalln(fmt.Errorf("couldn't create a file watch for file %s: %s", cfg.FilePath, err))
		}
		fmt.Printf("watching %s\n", cfg.FilePath)

		fileInfo, err := os.Stat(cfg.FilePath)
		if err != nil {
			log.Fatalln(fmt.Errorf("failed to get file info: %s", err))
		}

		if fileInfo.ModTime().After(time.Now().Add(-time.Minute * 5)) {
			OpenIDInGUDecks(cfg.FilePath)
		} else {
			fmt.Printf("file found, but older than 5min, not opening\n")
		}

		found := false
		for e := range watcher.Events {
			if e.Op == fsnotify.Write || e.Op == fsnotify.Create {
				if !found {
					err = OpenIDInGUDecks(cfg.FilePath)
					if err != nil {
						continue
					}
					found = true
				}
			} else if e.Op == fsnotify.Remove {
				fmt.Printf("file removed, trying to recreate watcher\n")
				found = false
				break
			} else {
				fmt.Printf("event: %s\n", e.Name)
			}
		}
	}

	fmt.Printf("exiting... - shouldn't occur\n")
}

func OpenIDInGUDecks(filePath string) error {
	id, err := GetID(filePath)
	if err != nil {
		if !errors.Is(err, ErrNotFound) {
			log.Fatalln(fmt.Errorf("fatal error reading the file: %s", err))
		}

		return ErrNotFound
	}

	fmt.Printf("found ID: %s, ourID: %s\n", id, ourID)
	browser.OpenURL(fmt.Sprintf(`https://gudecks.com/meta/player-stats?gameMode=13&userId=%s`, id))
	return nil
}

// GetID
func GetID(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalln(fmt.Errorf("couldn't open file: %s", err))
	}

	// file is usually < a few mb
	data, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	file.Close()

	cursor := 0

	for {
		userIDIndex := bytes.Index(data[cursor:], prefix)
		if userIDIndex == -1 {
			return "", ErrNotFound
		}

		userIDIndex += cursor + len(prefix)

		endIndex := bytes.IndexAny(data[userIDIndex:], " \n,)")
		if endIndex == -1 {
			return "", ErrNotFound
		}

		endIndex += userIDIndex
		cursor = endIndex
		userID := string(data[userIDIndex:endIndex])
		if userID == ourID {
			continue
		}

		if ourID == "" {
			ourID = userID
			continue
		} else {
			fmt.Printf("found after %d bytes\n", endIndex)
			return userID, nil
		}
	}
}
