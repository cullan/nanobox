package hookit

import (
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/display"
)

// Exec executes a hook inside of a container
func Exec(container, hook, payload, displayLevel string) (string, error) {
	return util.DockerExec(container, "root", "/opt/nanobox/hooks/"+hook, []string{payload}, display.NewStreamer(displayLevel))
}
