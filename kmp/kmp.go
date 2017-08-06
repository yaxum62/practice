package kmp

type Pattern struct {
	pattern string
	next    []int
}

func Compile(pattern string) *Pattern {
	if len(pattern) == 0 {
		return &Pattern{"", nil}
	}

	next := make([]int, len(pattern))
	next[0] = -1

	i := 0
	for i < len(pattern)-1 {
		j := next[i]
		for j >= 0 && pattern[j] != pattern[i] {
			j = next[j]
		}
		i++
		next[i] = j + 1
	}

	return &Pattern{
		pattern: pattern,
		next:    next,
	}
}

func (p *Pattern) FindIn(text string) int {
	if len(p.pattern) == 0 {
		return 0
	}

	i := 0
	j := 0
	for i < len(text) {
		for j >= 0 && text[i] != p.pattern[j] {
			j = p.next[j]
		}

		i++
		j++

		if j == len(p.pattern) {
			return i - j
		}
	}

	return -1
}
