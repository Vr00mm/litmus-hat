package types

type Hub struct {
	NAME    string
	URL     string
	BRANCH  string
	SSH_KEY string
}

type GitOPS struct {
	URL     string
	BRANCH  string
	SSH_KEY string
}

type Config struct {
	Projects []Project
}
type Project struct {
	Name   string
	Hubs   []Hub
	GitOPS GitOPS
}

type ClientLitmus struct {
	ACCESS_TOKEN string
	URL          string
}

type JWTLitmus struct {
	EXP      int    `json:exp`
	ROLE     string `json:role`
	UID      string `json:uid`
	USERNAME string `json:username`
}

type JWTToken struct {
	ACCESS_TOKEN string `json:token`
	EXPIRES_IN   int    `json:expires_in`
	TYPE         string `json:type`
}

type HubsList struct {
	Data struct {
		GetHubStatus []struct {
			ID            string `json:"id"`
			HubName       string `json:"HubName"`
			RepoBranch    string `json:"RepoBranch"`
			RepoURL       string `json:"RepoURL"`
			TotalExp      string `json:"TotalExp"`
			IsAvailable   bool   `json:"IsAvailable"`
			AuthType      string `json:"AuthType"`
			IsPrivate     bool   `json:"IsPrivate"`
			Token         string `json:"Token"`
			UserName      string `json:"UserName"`
			Password      string `json:"Password"`
			SSHPrivateKey string `json:"SSHPrivateKey"`
			SSHPublicKey  string `json:"SSHPublicKey"`
			LastSyncedAt  string `json:"LastSyncedAt"`
			Typename      string `json:"__typename"`
		} `json:"getHubStatus"`
	} `json:"data"`
}

type ProjectsList struct {
	Data []struct {
		ID      string `json:"ID"`
		UID     string `json:"UID"`
		Name    string `json:"Name"`
		Members []struct {
			UserID        string      `json:"UserID"`
			UserName      string      `json:"UserName"`
			Name          string      `json:"Name"`
			Role          string      `json:"Role"`
			Email         string      `json:"Email"`
			Invitation    string      `json:"Invitation"`
			JoinedAt      string      `json:"JoinedAt"`
			DeactivatedAt interface{} `json:"DeactivatedAt"`
		} `json:"Members"`
		State     string `json:"State"`
		CreatedAt string `json:"CreatedAt"`
		UpdatedAt string `json:"UpdatedAt"`
		RemovedAt string `json:"RemovedAt"`
	} `json:"data"`
}
