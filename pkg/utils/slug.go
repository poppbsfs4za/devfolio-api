package utils

import "strings"

func Slugify(input string) string {
	slug := strings.ToLower(strings.TrimSpace(input))
	replacer := strings.NewReplacer(
		" ", "-",
		"_", "-",
		"/", "-",
		"--", "-",
	)
	slug = replacer.Replace(slug)
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}
	return strings.Trim(slug, "-")
}
