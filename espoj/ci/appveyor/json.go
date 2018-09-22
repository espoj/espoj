package appveyor

import "strings"

func getAllFromJson(json, d string) []string {
	strs := make([]string, 0)
	si := 0
	var se int
	indexAddLen := len(d) + 3
	for si != -1 {
		preSi := si
		si = strings.Index(json[si:], d)
		if si == -1 {
			break
		}
		si += indexAddLen + preSi
		se = strings.Index(json[si:], "\"") + si
		strs = append(strs, json[si:se])
		si = se
	}
	return strs
}