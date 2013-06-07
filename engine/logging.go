package engine

import "github.com/jeady/go-logging"

var log *logging.Logger

func init() {
  log = logging.MustGetLogger("lmk")
}
