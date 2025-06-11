package models

/*example

curl https://git/rest/api/1.0/projects/FNX/repos/vessel/pull-requests/2436/activities

*/

type Activities struct {
	Size       int  `json:"size"`
	Limit      int  `json:"limit"`
	IsLastPage bool `json:"isLastPage"`
	Values     []struct {
		ID          int   `json:"id"`
		CreatedDate int64 `json:"createdDate"`
		User        struct {
			Name         string `json:"name"`
			EmailAddress string `json:"emailAddress"`
			Active       bool   `json:"active"`
			DisplayName  string `json:"displayName"`
			ID           int    `json:"id"`
			Slug         string `json:"slug"`
			Type         string `json:"type"`
			Links        struct {
				Self []struct {
					Href string `json:"href"`
				} `json:"self"`
			} `json:"links"`
		} `json:"user"`
		Action           string `json:"action"`
		FromHash         string `json:"fromHash,omitempty"`
		PreviousFromHash string `json:"previousFromHash,omitempty"`
		PreviousToHash   string `json:"previousToHash,omitempty"`
		ToHash           string `json:"toHash,omitempty"`
		Added            struct {
			Commits []struct {
				ID        string `json:"id"`
				DisplayID string `json:"displayId"`
				Author    struct {
					Name         string `json:"name"`
					EmailAddress string `json:"emailAddress"`
					Active       bool   `json:"active"`
					DisplayName  string `json:"displayName"`
					ID           int    `json:"id"`
					Slug         string `json:"slug"`
					Type         string `json:"type"`
					Links        struct {
						Self []struct {
							Href string `json:"href"`
						} `json:"self"`
					} `json:"links"`
				} `json:"author"`
				AuthorTimestamp int64 `json:"authorTimestamp"`
				Committer       struct {
					Name         string `json:"name"`
					EmailAddress string `json:"emailAddress"`
					Active       bool   `json:"active"`
					DisplayName  string `json:"displayName"`
					ID           int    `json:"id"`
					Slug         string `json:"slug"`
					Type         string `json:"type"`
					Links        struct {
						Self []struct {
							Href string `json:"href"`
						} `json:"self"`
					} `json:"links"`
				} `json:"committer"`
				CommitterTimestamp int64  `json:"committerTimestamp"`
				Message            string `json:"message"`
				Parents            []struct {
					ID        string `json:"id"`
					DisplayID string `json:"displayId"`
				} `json:"parents"`
				Properties struct {
					JiraKey []string `json:"jira-key"`
				} `json:"properties"`
			} `json:"commits"`
			Total int `json:"total"`
		} `json:"added,omitempty"`
		Removed struct {
			Commits []interface{} `json:"commits"`
			Total   int           `json:"total"`
		} `json:"removed,omitempty"`
		CancelledReason string  `json:"cancelledReason,omitempty"`
		CommentAction   string  `json:"commentAction,omitempty"`
		Comment         Comment `json:"comment,omitempty"`
		CommentAnchor   struct {
			FromHash string `json:"fromHash"`
			ToHash   string `json:"toHash"`
			Line     int    `json:"line"`
			LineType string `json:"lineType"`
			FileType string `json:"fileType"`
			Path     string `json:"path"`
			DiffType string `json:"diffType"`
			Orphaned bool   `json:"orphaned"`
		} `json:"commentAnchor,omitempty"`
		Diff struct {
			Source      interface{} `json:"source"`
			Destination struct {
				Components []string `json:"components"`
				Parent     string   `json:"parent"`
				Name       string   `json:"name"`
				Extension  string   `json:"extension"`
				ToString   string   `json:"toString"`
			} `json:"destination"`
			Hunks []struct {
				SourceLine      int `json:"sourceLine"`
				SourceSpan      int `json:"sourceSpan"`
				DestinationLine int `json:"destinationLine"`
				DestinationSpan int `json:"destinationSpan"`
				Segments        []struct {
					Type  string `json:"type"`
					Lines []struct {
						Destination int    `json:"destination"`
						Source      int    `json:"source"`
						Line        string `json:"line"`
						Repository  bool   `json:"repository"`
					} `json:"lines"`
					Truncated bool `json:"truncated"`
				} `json:"segments"`
				Truncated bool `json:"truncated"`
			} `json:"hunks"`
			Truncated  bool `json:"truncated"`
			Properties struct {
				ToHash   string `json:"toHash"`
				Current  bool   `json:"current"`
				FromHash string `json:"fromHash"`
			} `json:"properties"`
		} `json:"diff,omitempty"`
	} `json:"values"`
	Start int `json:"start"`
}
