package command

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Arguments struct {
	Repository string
	Commit     string
	OrderBy    string
	UseAuthor  bool
	Format     string
	Exclude    [][]string
	RestrictTo [][]string
	Extensions [][]string
}

type language struct {
	Name       string   `json:"name"`
	Type       string   `json:"type"`
	Extensions []string `json:"extensions"`
}

func findRoot() string {
	root, _ := os.Getwd()
	_, err := os.Stat(root + "/go.mod")
	for err != nil {
		root = filepath.Dir(root)
		_, err = os.Stat(root + "/go.mod")
	}
	return root
}

func loadLanguages() map[string][]string {
	var languageJSON []language
	languageDict := make(map[string][]string)
	fileBytes, _ := os.ReadFile(
		findRoot() + "/gitfame/configs/language_extensions." +
			"json",
	)
	_ = json.Unmarshal(fileBytes, &languageJSON)
	for _, language := range languageJSON {
		languageDict[strings.ToLower(language.Name)] = language.Extensions
	}
	return languageDict
}

func ParseArguments(args []string) Arguments {
	log.Println("parsing arguments")
	languageDict := loadLanguages()
	arguments := Arguments{
		".", "HEAD", "lines", true, "tabular", [][]string{},
		[][]string{}, [][]string{},
	}
	arguments.UseAuthor = true
	for index := 1; index < len(args); index++ {
		command := args[index]
		splitIndex := strings.Index(command, "=")
		var value string
		if splitIndex != -1 {
			value = command[splitIndex+1:]
			command = command[:splitIndex]
		}
		switch command {
		case "--repository":
			if splitIndex == -1 {
				index++
				value = args[index]
			}
			arguments.Repository = value

		case "--revision":
			if splitIndex == -1 {
				index++
				value = args[index]
			}
			arguments.Commit = value
		case "--order-by":
			if splitIndex == -1 {
				index++
				value = args[index]
			}
			arguments.OrderBy = value
		case "--use-committer":
			arguments.UseAuthor = false
		case "--format":
			if splitIndex == -1 {
				index++
				value = args[index]
			}
			arguments.Format = value
		case "--extensions":
			if splitIndex == -1 {
				index++
				value = args[index]
			}
			extensions := strings.Split(value, ",")
			arguments.Extensions = append(arguments.Extensions, extensions)
		case "--languages":
			if splitIndex == -1 {
				index++
				value = args[index]
			}
			languages := strings.Split(value, ",")
			var allExtensions []string
			for _, language := range languages {
				extensions, contains := languageDict[strings.ToLower(language)]
				if contains {
					allExtensions = append(allExtensions, extensions...)
				}
			}
			arguments.Extensions = append(
				arguments.Extensions, allExtensions,
			)
		case "--restrict-to":
			if splitIndex == -1 {
				index++
				value = args[index]
			}
			arguments.RestrictTo = append(
				arguments.RestrictTo, strings.Split(value, ","),
			)
		case "--exclude":
			if splitIndex == -1 {
				index++
				value = args[index]
			}
			arguments.Exclude = append(
				arguments.Exclude, strings.Split(value, ","),
			)

		}
	}
	return arguments
}
