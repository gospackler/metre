package metre

import (
	"strings"

	"go.uber.org/zap"

	"github.com/gospackler/metre/logging"
	"github.com/gospackler/metre/transport"
)

func StartBroker(dealerUri string, routerUri string) {
	if routerUri == "" {
		routerUri = LocalHost + ":" + RouterPort
	} else if strings.Index(routerUri, ":") == 0 {
		routerUri = LocalHost + ":" + RouterPort
	}

	if dealerUri == "" {
		dealerUri = LocalHost + ":" + DealerPort
	} else if strings.Index(dealerUri, ":") == 0 {
		dealerUri = LocalHost + ":" + DealerPort
	}

	go func() {
		err := transport.StartBroker(dealerUri, routerUri)
		if err != nil {
			logging.Logger.Warn("error starting broker",
				zap.Error(err),
			)
		}
	}()
}
