package models

/*example
TODO implement this
curl -su user https://git/rest/api/1.0/projects/FNX/repos/bank-transactions-fetcher/pull-requests/136/merge

{
  "canMerge": false,
  "conflicted": false,
  "outcome": "CLEAN/CONFLICTED",
  "vetoes": [
    {
      "summaryMessage": "Not all required reviewers have approved yet",
      "detailedMessage": "or must review and approve this pull request before it can be merged."
    }
  ]
}

*/

//
type Merge struct {
	CanMerge   bool   `json:"canMerge"`
	Conflicted bool   `json:"conflicted"`
	Outcome    string `json:"outcome"`
	Vetoes     []Veto `json:"vetoes"`
}

type Veto struct {
	SummaryMessage  string `json:"summaryMessage"`
	DetailedMessage string `json:"detailedMessage"`
}
