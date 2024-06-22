package environment

import (
	"errors"
	"log"
	"os"
	"strings"
)

var path string

func init() {
	path, _ = os.Getwd()
}

func ValidateCreation(flags map[string]bool) []error {

	var errs []error
	if res, err := os.Stat(".envim"); err == nil {
		if res.IsDir() {
			errs = append(errs, errors.New("Environment folder already exists in "+path))
		} else {
			errs = append(errs, errors.New("A file by the name '.envim' already exists in the current directory"))
		}
	}

	for key, val := range flags {
		if val {
			switch key {
			case "dotnvim":
				if res, err := os.Stat(".nvim"); err == nil {
					if res.IsDir() {
						errs = append(errs, errors.New("'.nvim' folder already exists in "+path))
					} else {
						errs = append(errs, errors.New("A file by the name '.nvim' already exists in the current directory"))
					}
				}
			case "gitignore":
				if res, err := os.Stat(".gitignore"); err == nil && res.IsDir() {
					errs = append(errs, errors.New("There is a folder named '.gitignore' in the current directory"))
				}
			}
		}
	}

	return errs
}

func CreateEnvironment() error {
	if err := os.Mkdir(".envim", 0755); err != nil {
		return err
	}
	log.Printf("Environment created in %s\n", path)
	return nil
}

func CreateDotNvim() error {
	if err := os.Mkdir(".nvim", 0755); err != nil {
		return err
	}
	log.Printf("Dotnvim folder created in %s\n", path)
	return nil
}

func AppendToGitignore() error {

	dat, err := os.ReadFile(".gitignore")
	if err == nil {
		for _, line := range strings.Split(string(dat), "\n") {
			if line == ".envim" || line == ".envim/*" {
				log.Printf(".envim already appended to .gitignore in %s\n", path)
				return nil
			}
		}
	}

	file, err := os.OpenFile(".gitignore", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.WriteString(".envim\n"); err != nil {
		return err
	}
	log.Printf(".envim appended to .gitignore in %s\n", path)
	return nil
}
