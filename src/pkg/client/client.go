package client

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"litmus-init/pkg/tools"
	"litmus-init/types"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

type clientLitmus struct {
	ACCESS_TOKEN string
	URL          string
}

func NewClient(url string, username string, password string) clientLitmus {
	var client clientLitmus
	var displayPass string

	if log.IsLevelEnabled(log.DebugLevel) {
		displayPass = password
	} else {
		displayPass = "***REDACTED***"
	}
	log.WithFields(log.Fields{"LITMUS_URL": url, "LITMUS_ADMIN": username, "LITMUS_PASSWORD": displayPass}).Info("Initialisation")

	req := `{ "username": "` + username + `", "password": "` + password + `"}`
	resp, err := http.Post(url+"/auth/login", "Application/json", bytes.NewReader([]byte(req)))
	if err != nil {
		log.Fatalf("Fatal error: cannot login\n%s\n", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Fatal error: cannot parse response\n%s\n", err)
	}

	var rawToken types.JWTToken
	err = json.Unmarshal(body, &rawToken)
	if err != nil {
		log.Fatalf("Fatal error: cannot unmarshal response\n%s\n##########################\n%s\n", err, body)
	}

	if rawToken.ACCESS_TOKEN == "" {
		log.Fatalf("Fatal error: cannot login to litmus, err:\n%s\n##########################\n", body)
	}

	client.ACCESS_TOKEN = rawToken.ACCESS_TOKEN
	client.URL = url
	return client
}

func (client clientLitmus) GetUID() string {
	tmp := strings.Split(client.ACCESS_TOKEN, ".")

	dataLen := len(tmp[1]) % 4
	var result string

	switch dataLen {
	case 2:
		result = tmp[1] + "=="
	case 3:
		result = tmp[1] + "="
	default:
		result = tmp[1]
	}

	decoded, err := base64.StdEncoding.DecodeString(result)
	if err != nil {
		log.Fatalf("Fatal error: cannot base64 decode JWT\n%s\n", err)
	}

	var litmusToken types.JWTLitmus
	err = json.Unmarshal(decoded, &litmusToken)
	if err != nil {
		log.Fatalf("Fatal error: cannot unmarshall JWT\n%s\n", err)
	}

	return litmusToken.UID

}

func (client clientLitmus) CreateProject(projectName string) {
	var createProject tools.Request
	createProject.METHOD = "POST"
	createProject.ENDPOINT = client.URL + "/auth/create_project"
	createProject.REQUEST = `{"project_name":"` + projectName + `", "user_id":"` + string(client.GetUID()) + `"}`

	_ = tools.GraphqlRequest(createProject, client.ACCESS_TOKEN)
}

func (client clientLitmus) CreateHub(hub types.Hub, projectName string) {
	var createHub tools.Request
	createHub.METHOD = "POST"
	createHub.ENDPOINT = client.URL + "/api/query"

	hubSshKey := strings.Replace(hub.SSH_KEY, "\n", "\x5Cn", -1)

	createHub.REQUEST = `{"operationName":"addMyHub","variables":{"MyHubDetails":{"HubName":"` + hub.NAME + `","RepoURL":"` + hub.URL + `","RepoBranch":"` + hub.BRANCH + `","IsPrivate":true,"AuthType":"ssh","Token":"","UserName":"user","Password":"user","SSHPrivateKey":"` + hubSshKey + `","SSHPublicKey":"TO_SET"},"projectID":"` + client.GetProjectID(projectName) + `"},"query":"mutation addMyHub($MyHubDetails: CreateMyHub!, $projectID: String!) {\n  addMyHub(myhubInput: $MyHubDetails, projectID: $projectID) {\n    HubName\n    RepoURL\n    RepoBranch\n    __typename\n  }\n}\n"}`

	_ = tools.GraphqlRequest(createHub, client.ACCESS_TOKEN)
}

func (client clientLitmus) GetProjectID(projectName string) string {
	var projectID string

	var getProjectID tools.Request
	getProjectID.METHOD = "GET"
	getProjectID.ENDPOINT = client.URL + "/auth/list_projects"
	getProjectID.REQUEST = ``

	result := tools.GraphqlRequest(getProjectID, client.ACCESS_TOKEN)

	var projectsList types.ProjectsList
	err := json.Unmarshal([]byte(result), &projectsList)
	if err != nil {
		log.Fatalf("Fatal error: cannot unmarshall projectlist\n%s\n", err)
	}

	for _, project := range projectsList.Data {
		if project.Name == projectName {
			projectID = project.ID
		}

	}
	return projectID
}

func (client clientLitmus) DeleteHub(hubID string, projectName string) {
	var deleteHub tools.Request
	deleteHub.METHOD = "POST"
	deleteHub.ENDPOINT = client.URL + "/api/query"
	deleteHub.REQUEST = `{"operationName":"deleteMyHub","variables":{"hub_id":"` + hubID + `","projectID":"` + client.GetProjectID(projectName) + `"},"query":"mutation deleteMyHub($hub_id: String!, $projectID: String!) {\n  deleteMyHub(hub_id: $hub_id, projectID: $projectID)\n}\n"}`
	_ = tools.GraphqlRequest(deleteHub, client.ACCESS_TOKEN)

}

func (client clientLitmus) GetHubsList(projectName string) types.HubsList {
	var getHubsList tools.Request
	getHubsList.METHOD = "POST"
	getHubsList.ENDPOINT = client.URL + "/api/query"
	getHubsList.REQUEST = `{"operationName":"getHubStatus","variables":{"data":"` + client.GetProjectID(projectName) + `"},"query":"query getHubStatus($data: String! ) {\n  getHubStatus(projectID: $data) {\n    id\n    HubName\n    RepoBranch\n    RepoURL\n    TotalExp\n    IsAvailable\n    AuthType\n    IsPrivate\n    Token\n    UserName\n    Password\n    SSHPrivateKey\n    SSHPublicKey\n    LastSyncedAt\n    __typename\n  }\n}\n"}`

	result := tools.GraphqlRequest(getHubsList, client.ACCESS_TOKEN)

	var hubsList types.HubsList
	err := json.Unmarshal([]byte(result), &hubsList)
	if err != nil {
		log.Fatalf("Fatal error: cannot unmarshall projectlist\n%s\n", err)
	}

	return hubsList
}

func (client clientLitmus) DeleteAllHubs(projectName string) {
	hubsList := client.GetHubsList(projectName)
	for _, hub := range hubsList.Data.GetHubStatus {
		client.DeleteHub(hub.ID, projectName)
	}
}

func (client clientLitmus) ConfigureGitOPS(gitops types.GitOPS, projectName string) {
	gitSshKey := strings.Replace(gitops.SSH_KEY, "\n", "\x5Cn", -1)

	var configureGitOPS tools.Request
	configureGitOPS.METHOD = "POST"
	configureGitOPS.ENDPOINT = client.URL + "/api/query"
	configureGitOPS.REQUEST = `{"operationName":"enableGitOps","variables":{"gitConfig":{"ProjectID":"` + client.GetProjectID(projectName) + `","RepoURL":"` + gitops.URL + `","Branch":"` + gitops.BRANCH + `","AuthType":"ssh","Token":"","UserName":"user","Password":"user","SSHPrivateKey":"` + gitSshKey + `"}},"query":"mutation enableGitOps($gitConfig: GitConfig!) {\n  enableGitOps(config: $gitConfig)\n}\n"}`

	_ = tools.GraphqlRequest(configureGitOPS, client.ACCESS_TOKEN)
}
