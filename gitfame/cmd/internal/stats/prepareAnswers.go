package stats

import (
	"gitlab.com/slon/shad-go/gitfame/cmd/internal/command"
	"gitlab.com/slon/shad-go/gitfame/cmd/internal/gitcommand"
	"sort"
)

type FinalAuthorInfo struct {
	Author  string
	Files   int
	Commits int
	Lines   int
}

func defaultComparator(
	info1 *FinalAuthorInfo,
	info2 *FinalAuthorInfo,
) bool {
	return info1.Lines > info2.Lines || info1.Lines == info2.Lines && (info1.
		Commits > info2.Commits || info1.Commits == info2.Commits && (info1.
		Files > info2.Files || info1.Files == info2.Files && info1.
		Author < info2.Author))
}

func CreateAnswer(
	intermediateAuthorInfo map[string]*gitcommand.AuthorInfo,
	args command.Arguments,
) []FinalAuthorInfo {
	var authorInfo []FinalAuthorInfo
	for k, v := range intermediateAuthorInfo {
		authorInfo = append(
			authorInfo, FinalAuthorInfo{
				Author: k,
				Files:  len(v.Files), Commits: len(v.Commits),
				Lines: v.LinesNumber,
			},
		)
	}
	sort.Slice(
		authorInfo, func(i, j int) bool {
			if args.OrderBy == "commits" {
				return authorInfo[i].Commits > authorInfo[j].Commits ||
					authorInfo[i].Commits == authorInfo[j].
						Commits && defaultComparator(
						&authorInfo[i], &authorInfo[j],
					)
			} else if args.OrderBy == "files" {
				return authorInfo[i].Files > authorInfo[j].Files ||
					authorInfo[i].Files == authorInfo[j].
						Files && defaultComparator(
						&authorInfo[i], &authorInfo[j],
					)
			} else if args.OrderBy == "lines" {
				return defaultComparator(&authorInfo[i], &authorInfo[j])
			} else {
				panic(1)
			}
		},
	)
	return authorInfo
}
