package docopt

type ArgumentsMatcher struct{}

func (matcher *ArgumentsMatcher) Match(
	args []string,
	variants []Grammar,
	options []Option,
) (map[string]interface{}, error) {
	//var result map[string]interface{}

	//for _, variant := range variants {
	//    for {
	//        for _, token := range variant {
	//            matcher.matchToken(args, token, result)
	//        }
	//    }
	//}

	//return result, nil
	return nil, nil
}

//func (matcher *ArgumentsMatcher) matchToken(
//    arg string,
//    token Token,
//    result map[string]interface{},
//) (tail string, matched bool) {
//    switch token := token.(type) {
//    case TokenSeparator:
//        if arg != "" {
//            return "", false
//        }

//        return "", true

//    case TokenOption:
//        return "", false
//    }
//}
