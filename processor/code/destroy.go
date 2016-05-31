package code

import (
	"net"

	"github.com/nanobox-io/golang-docker-client"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/provider"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/data"
	"github.com/nanobox-io/nanobox/util/ip_control"
)

type codeDestroy struct {
	control processor.ProcessControl
}

func init() {
	processor.Register("code_destroy", codeDestroyFunc)
}

func codeDestroyFunc(control processor.ProcessControl) (processor.Processor, error) {
	// confirm the provider is an accessable one that we support.
	if control.Meta["name"] == "" {
		return nil, missingImageOrName
	}
	return &codeDestroy{control: control}, nil
}

func (self codeDestroy) Results() processor.ProcessControl {
	return self.control
}

func (self *codeDestroy) Process() error {

	// get the service from the database
	service := models.Service{}
	err := data.Get(util.AppName(), self.control.Meta["name"], &service)
	if err != nil {
		// cant find service
		return err
	}

	err = docker.ContainerRemove(service.ID)
	if err != nil {
		return err
	}

	err = provider.RemoveNat(service.ExternalIP, service.InternalIP)
	if err != nil {
		return err
	}

	err = provider.RemoveIP(service.ExternalIP)
	if err != nil {
		return err
	}

	err = ip_control.ReturnIP(net.ParseIP(service.ExternalIP))
	if err != nil {
		return err
	}

	err = ip_control.ReturnIP(net.ParseIP(service.InternalIP))
	if err != nil {
		return err
	}

	// save the service
	return data.Delete(util.AppName(), self.control.Meta["name"])
}
