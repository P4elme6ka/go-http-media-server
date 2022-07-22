package app

import (
	"fmt"
	"github.com/P4elme6ka/go-http-media-server/src/goVirtualHost"
	"github.com/P4elme6ka/go-http-media-server/src/util"
)

func printAccessibleURLs(vhSvc *goVirtualHost.Service) {
	vhostsUrls := vhSvc.GetAccessibleURLs(false)
	file, teardown := util.GetTTYFile()

	for vhIndex := range vhostsUrls {
		fmt.Fprintln(file, "Host", vhIndex, "may be accessed by URLs:")
		for urlIndex := range vhostsUrls[vhIndex] {
			fmt.Fprintln(file, "  ", vhostsUrls[vhIndex][urlIndex])
		}
	}

	teardown()
}
