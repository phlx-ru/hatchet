package strings

func Unique(in []string) []string {
	hashmap := map[string]struct{}{}

	out := []string{}
	for _, item := range in {
		if _, ok := hashmap[item]; ok {
			continue
		}
		hashmap[item] = struct{}{}
		out = append(out, item)
	}

	return out
}
