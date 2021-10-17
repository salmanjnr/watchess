package broadcast

import (
	"fmt"

	"github.com/salman69e27/chess"
)

func encodePGNWithComments(g *chess.Game) string {
	s := ""
	tagPairs := g.TagPairs()
	moves := g.Moves()
	positions := g.Positions()
	comments := g.Comments()
	notation := chess.AlgebraicNotation{}
	for _, tag := range tagPairs {
		s += fmt.Sprintf("[%s \"%s\"]\n", tag.Key, tag.Value)
	}
	s += "\n"
	for i, move := range moves {
		pos := positions[i]
		txt := notation.Encode(pos, move)
		if i%2 == 0 {
			if i == 0 {
				s += fmt.Sprintf("%d. %s", (i/2)+1, txt)
			} else {
				s += fmt.Sprintf(" %d. %s", (i/2)+1, txt)
			}
			for _, comment := range comments[i] {
				s += fmt.Sprintf(" { %s }", comment)
			}
		} else {
			if len(comments[i-1]) == 0 {
				s += fmt.Sprintf(" %s", txt)
			} else {
				s += fmt.Sprintf(" %d... %s", ((i-1)/2)+1, txt)
				for _, comment := range comments[i] {
					s += fmt.Sprintf(" { %s }", comment)
				}
			}
		}
	}
	s += " " + g.Outcome().String()
	return s
}
