package alertmanager

type AlertmanagerSilence struct {
	Matchers []struct {
		Name    string `json:"name"`
		Value   string `json:"value"`
		IsRegex bool   `json:"isRegex"`
	} `json:"matchers"`
	StartsAt  string `json:"startsAt"`
	EndsAt    string `json:"endsAt"`
	CreatedBy string `json:"createdBy"`
	Comment   string `json:"comment"`
}

type AlertManagerFederation struct {
	Env                    string            `json:"env"`
	ParentAlertManagerURL  string            `json:"parentAlertManagerURL"`
	ChildsAlertManagerURLs map[string]string `json:"childsAlertManagerURLs"`
}
