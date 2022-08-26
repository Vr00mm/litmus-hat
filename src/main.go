package main

import (
	clientLitmus "litmus-init/pkg/client"
	"litmus-init/pkg/config"
	clientK8s "litmus-init/pkg/k8s"
	"litmus-init/pkg/tools"
	"os"

	log "github.com/sirupsen/logrus"
)

var (
	LITMUS_URL      string
	LITMUS_ADMIN    string
	LITMUS_PASSWORD string
)

func init() {
	LITMUS_URL = os.Getenv("LITMUS_URL")
	LITMUS_ADMIN = os.Getenv("ADMIN_USERNAME")
	LITMUS_PASSWORD = os.Getenv("ADMIN_PASSWORD")
}

func main() {
	tools.SetLogLevel()

	log.Info("Waiting for litmus components to be ready")

	clientK8s.WaitForPodBySelectorRunning("litmus", "helm.sh/chart=litmus-2.8.2", 300)

	client := clientLitmus.NewClient(LITMUS_URL, LITMUS_ADMIN, LITMUS_PASSWORD)

	configData := config.Load()
	for _, projectData := range configData.Projects {

		log.Infoln("Create Project: " + projectData.Name)
		client.CreateProject(projectData.Name)

		log.Infoln("Setup Hubs for: " + projectData.Name)
		client.DeleteAllHubs(projectData.Name)

		for _, hubData := range projectData.Hubs {
			log.Infoln("Configure Hub: " + hubData.NAME)
			client.CreateHub(hubData, projectData.Name)
		}

		log.Infoln("Setup GITOPS for: " + projectData.Name)
		client.ConfigureGitOPS(projectData.GitOPS, projectData.Name)

	}
	os.Exit(0)

}
