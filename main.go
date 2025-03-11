package main

import (
	"runtime/debug"
	"time"

	"github.com/mheers/raygun2x/cmd"
	"github.com/mheers/raygun2x/models"
	"github.com/sirupsen/logrus"
)

// const raygunFolder = "presets/"

// func main() {
// 	suites, err := raygun.GetTestSuites(raygunFolder)
// 	if err != nil {
// 		panic(err)
// 	}

// 	postmanCollection, err := postman.CreateCollection(suites)
// 	if err != nil {
// 		panic(err)
// 	}
// 	if err := postman.WriteToFile(postmanCollection, "raygun.postman_collection.json"); err != nil {
// 		panic(err)
// 	}

// 	openAPISpec, err := openapi.CreateOpenAPISpec(suites)
// 	if err != nil {
// 		panic(err)
// 	}
// 	if err := openapi.WriteToFile(openAPISpec, "raygun.openapi.yaml"); err != nil {
// 		panic(err)
// 	}
// }

// build flags
var (
	VERSION   string
	BuildTime string
	GitTag    string
	GitBranch string
)

func main() {
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		panic("could not read build info")
	}

	runtimeInfo := &models.RuntimeInfo{}

	for _, kv := range buildInfo.Settings {
		switch kv.Key {
		case "vcs.revision":
			runtimeInfo.CommitHash = kv.Value
		case "vcs.time":
			runtimeInfo.CommitTime, _ = time.Parse(time.RFC3339, kv.Value)
		case "vcs.modified":
			runtimeInfo.DirtyBuild = kv.Value == "true"
		}
	}

	runtimeInfo.Version = VERSION
	runtimeInfo.BuildTime = BuildTime
	runtimeInfo.GoVersion = buildInfo.GoVersion
	runtimeInfo.GitTag = GitTag
	runtimeInfo.GitBranch = GitBranch

	cmd.RuntimeInfo = runtimeInfo

	// execeute the command
	err := cmd.Execute()
	if err != nil {
		logrus.Fatalf("Execute failed: %+v", err)
	}
}
