package api

import (
	"github.com/stretchr/testify/require"
	db "github.com/yerakairzhan/SimpleBank/db/sqlc"
	"github.com/yerakairzhan/SimpleBank/util"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	config := util.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(config, store)
	require.NoError(t, err)

	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
