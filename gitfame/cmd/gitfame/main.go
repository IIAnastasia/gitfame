//go:build !solution

package main

import (
	"gitlab.com/slon/shad-go/gitfame/cmd/internal/command"
	"gitlab.com/slon/shad-go/gitfame/cmd/internal/gitcommand"
	"gitlab.com/slon/shad-go/gitfame/cmd/internal/printer"
	"gitlab.com/slon/shad-go/gitfame/cmd/internal/stats"
	"os"
)

func main() {
	arguments := command.ParseArguments(os.Args)
	files := gitcommand.FindFiles(arguments)
	authorsInfo := make(map[string]*gitcommand.AuthorInfo)
	for _, file := range files {
		gitcommand.Analyze(file, authorsInfo, arguments)
	}
	answer := stats.CreateAnswer(authorsInfo, arguments)
	printer.PrintInfo(arguments, answer)

}
