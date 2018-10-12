package main

type WebHookStruct struct {
	Url         string
	Events      []string
	Enabled     string
	SSLCheck    bool
	identified  bool
	ContentType string //json or from. default form
	name        string
}

const hookIgnored = "ignore"
