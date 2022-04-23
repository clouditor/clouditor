package evidence

// SetPageToken implements PaginatedRequest so we can set the page token programmatically.
func (r *ListEvidencesRequest) SetPageToken(token string) {
	r.PageToken = token
}
