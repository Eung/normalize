package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/text/unicode/norm"
)

var recursive bool

func main() {
	flag.BoolVar(&recursive, "r", false, "is recursive")
	flag.Parse()

	// fmt.Printf("recursive: %v\n", recursive)
	for idx, arg := range flag.Args() {
		fmt.Printf("arg[%d]: %s\n", idx, arg)
	}

	workdir := make([]string, 0)
	if flag.NArg() == 0 {
		workdir = append(workdir, ".")
	} else {
		workdir = append(workdir, flag.Args()...)
	}

	for _, wd := range workdir {
		if err := walk(wd); err != nil {
			fmt.Printf("walk err: %s", err.Error())
			os.Exit(-1)
		}
	}
}

func walk(workdir string) error {
	file, err := os.Open(workdir)
	if err != nil {
		fmt.Printf("[%s] open err: %s\n", workdir, err.Error())
		return err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		fmt.Printf("[%s] stat err: %s\n", workdir, err.Error())
		return err
	}

	if stat.IsDir() {
		fileinfo, err := file.Readdir(0)
		if err != nil {
			fmt.Printf("[%s] readdir err: %s\n", workdir, err.Error())
			return err
		}

		for _, file := range fileinfo {
			if file.Mode().IsDir() {
				if recursive {
					if err := walk(filepath.Join(workdir, file.Name())); err != nil {
						fmt.Printf("recursive err: %s\n", err.Error())
					}
				}
			}

			// fmt.Printf("visited: %s\n", filepath.Join(workdir, file.Name()))
			err := normalize(filepath.Join(workdir, file.Name()))
			if err != nil {
				fmt.Printf("normalize err: %s", err.Error())
				return err
			}
		}
	} else {
		// fmt.Printf("visited file: %s\n", file.Name())
		err := normalize(filepath.Join(file.Name()))
		if err != nil {
			fmt.Printf("normalize err: %s", err.Error())
			return err
		}
	}

	return nil
}

func normalize(oldPath string) error {
	newPath := filepath.Join(filepath.Dir(oldPath), norm.NFC.String(filepath.Base(oldPath)))

	if strings.Compare(oldPath, newPath) == 0 {
		return nil
	} else {
		return os.Rename(oldPath, newPath)
	}
}
