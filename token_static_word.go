package docopt

type TokenStaticWord struct {
	Name string
}

func (word *TokenStaticWord) String() string {
	return word.Name
}
