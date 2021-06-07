package metric_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"testing"

	"clouditor.io/clouditor/api/orchestrator"
	"clouditor.io/clouditor/cli"
	"clouditor.io/clouditor/cli/commands/login"
	"clouditor.io/clouditor/cli/commands/metric"
	service_auth "clouditor.io/clouditor/service/auth"
	service_orchestrator "clouditor.io/clouditor/service/orchestrator"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
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
	var b bytes.Buffer

	cli.Output = &b

	cmd := metric.NewListMetricsCommand()
	err = cmd.RunE(nil, []string{})

	assert.Nil(t, err)

	var response *orchestrator.ListMetricsResponse = &orchestrator.ListMetricsResponse{}

	err = protojson.Unmarshal(b.Bytes(), response)

	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.NotEmpty(t, response.Metrics)
}
