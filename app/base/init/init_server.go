package init

import (
	"fmt"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	basecommon "github.com/lyonmu/quebec/app/base/common"
	"github.com/sirupsen/logrus"
)

func Start() {

	gin.DefaultWriter = basecommon.System.Base.Log.GetWriter()
	gin.DefaultErrorWriter = basecommon.System.Base.Log.GetWriter()
	if basecommon.System.Base.Log.Level == "info" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(cors.Default())

	quebecBase := r.Group(basecommon.System.Base.RouterPrefix)
	RouterSetup(quebecBase)

	if err := r.Run(fmt.Sprintf(":%d", basecommon.System.Base.Port)); err != nil {
		logrus.Error("server start failed, err: ", err)
		os.Exit(1)
	}

}
