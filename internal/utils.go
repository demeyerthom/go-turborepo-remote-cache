package internal

func GetSlug(teamId, slug *string) *string {
	if teamId != nil {
		return teamId
	}

	return slug
}
