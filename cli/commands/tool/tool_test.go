package tool_test

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"testing"

	"clouditor.io/clouditor/api/orchestrator"
	"clouditor.io/clouditor/cli/commands/login"
	"clouditor.io/clouditor/cli/commands/tool"
	service_auth "clouditor.io/clouditor/service/auth"
	service_orchestrator "clouditor.io/clouditor/service/orchestrator"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

var sock net.Listener
var server *grpc.Server

func TestMain(m *testing.M) {
	var (
		err error
		dir string
	)

	sock, server, err = service_auth.StartDedicatedAuthServer(":0")
	orchestrator.RegisterOrchestratorServer(server, &service_orchestrator.Service{})

	if err != nil {
		panic(err)
	}

	defer sock.Close()
	defer server.Stop()

	dir, err = ioutil.TempDir(os.TempDir(), ".clouditor")
	if err != nil {
		panic(err)
	}

	viper.Set("username", "clouditor")
	viper.Set("password", "clouditor")
	viper.Set("session-directory", dir)

	cmd := login.NewLoginCommand()
	err = cmd.RunE(nil, []string{fmt.Sprintf("localhost:%d", sock.Addr().(*net.TCPAddr).Port)})
	if err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

func TestListTool(t *testing.T) {
	var err error

	cmd := tool.NewListToolsCommand()
	err = cmd.RunE(nil, []string{})

	// unsupported for now
	assert.NotNil(t, err)
	assert.Error(t, err, "method ListAssessmentTools not implemented")
}

func TestShowTool(t *testing.T) {
	var err error

	cmd := tool.NewShowToolCommand()
	err = cmd.RunE(nil, []string{"1"})

	// unsupported for now
	assert.NotNil(t, err)
	assert.Error(t, err, "method GetAssessmentTool not implemented")
}
