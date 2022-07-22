package clitest

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"testing"

	"clouditor.io/clouditor/api"
	"clouditor.io/clouditor/cli"
	"clouditor.io/clouditor/internal/testutil"
	"clouditor.io/clouditor/service"
	oauth2 "github.com/oxisto/oauth2go"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

// PrepareSession prepares a session for unit tests. It creates a temporary folder to save
// the session credentials in and does a login to the specified authorization server using
// test credentials. It is the responsibility of the caller to cleanup the temporary directory.
//
// This function will use asserts and fail/panic if errors occurs.
func PrepareSession(authPort uint16, authSrv *oauth2.AuthorizationServer, grpcURL string) (dir string, err error) {
	var (
		token   *oauth2.Token
		session *cli.Session
	)

	// Create a temporary folder
	dir, err = ioutil.TempDir(os.TempDir(), ".clouditor")
	if err != nil {
		return "", err
	}

	viper.Set("auth-server", fmt.Sprintf("http://localhost:%d", authPort))
	viper.Set(cli.SessionFolderFlag, dir)

	// Simulate a login by directly granting a token
	token, err = authSrv.GenerateToken(testutil.TestAuthClientID, 0, 0)
	if err != nil {
		return "", err
	}

	// TODO(oxisto): This is slightly duplicated code from the Login command. Extract it into the session struct
	session, err = cli.NewSession(grpcURL, &oauth2.Config{
		ClientID: testutil.TestAuthClientID,
		Endpoint: oauth2.Endpoint{
			AuthURL:  testutil.AuthURL(authPort),
			TokenURL: testutil.TokenURL(authPort),
		},
	}, token)
	if err != nil {
		return "", err
	}

	session.SetAuthorizer(api.NewOAuthAuthorizerFromConfig(
		session.Config,
		token,
	))

	err = session.Save()
	if err != nil {
		return "", err
	}

	return dir, nil
}

// RunCLITest can be used in a TestMain function for CLI tests. It takes care of launching
// an authorization server as well as a gRPC server with the selected services supplied as options.
// It also automatically issues a login command to the auth service.
//
// Since this function is primarily used in a TestMain and no testing.T object is available at this
// point, this function WILL panic on errors.
func RunCLITest(m *testing.M, opts ...service.StartGRPCServerOption) (code int) {
	var (
		err      error
		tmpDir   string
		auth     *oauth2.AuthorizationServer
		authPort uint16
		grpcPort uint16
		sock     net.Listener
		server   *grpc.Server
	)

	auth, authPort, err = testutil.StartAuthenticationServer()
	if err != nil {
		panic(err)
	}

	sock, server, err = service.StartGRPCServer(testutil.JWKSURL(authPort), opts...)
	if err != nil {
		panic(err)
	}

	grpcPort = sock.Addr().(*net.TCPAddr).AddrPort().Port()

	tmpDir, err = PrepareSession(authPort, auth, fmt.Sprintf("localhost:%d", grpcPort))
	if err != nil {
		panic(err)
	}

	code = m.Run()

	sock.Close()
	server.Stop()

	// Remove temporary session directory
	os.RemoveAll(tmpDir)

	return code
}

// AutoChdir automatically guesses if we need to change the current working directory
// so that we can find the policies folder
func AutoChdir() {
	var (
		err error
	)

	// Check, if we can find the policies folder
	_, err = os.Stat("policies")
	if errors.Is(err, os.ErrNotExist) {
		// Try again one level deeper
		err = os.Chdir("..")
		if err != nil {
			panic(err)
		}

		AutoChdir()
	}
}
