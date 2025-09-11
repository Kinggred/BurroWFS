package schemas

import (
	"burrowfs/api/schemas/webdav"
	"burrowfs/core/logging"
	"net/http"
	"strconv"

	"github.com/pgaskin/xmlwriter"
)

// rootCollectionObjectWriter writes the XML response for the root collection object. Required for WebDAV clients.
func rootCollectionObjectWriter(xw *xmlwriter.XMLWriter, davNS *xmlwriter.NS) error {
	err := xw.Start(davNS, "response")
	err = xw.Start(davNS, "href")
	err = xw.Text(false, "/")
	err = xw.End(false)
	err = xw.Start(davNS, "propstat")
	err = xw.Start(davNS, "prop")
	err = xw.Start(davNS, "resourcetype")
	err = xw.Start(davNS, "collection")
	err = xw.End(true)
	err = xw.End(false)
	err = xw.Start(davNS, "displayname")
	err = xw.Text(false, "/")
	err = xw.End(false)
	err = xw.End(false)
	err = xw.Start(davNS, "status")
	err = xw.Text(false, "HTTP/1.1 200 OK")
	err = xw.End(false)
	err = xw.End(false)
	err = xw.End(false)

	return err
}

func resourceObjectWriter(xw *xmlwriter.XMLWriter, davNS *xmlwriter.NS, resource webdav.ResourceSchema) error {
	err := xw.Start(davNS, "response")

	err = xw.Start(davNS, "href")
	err = xw.Text(false, resource.Href)
	err = xw.End(false)

	err = xw.Start(davNS, "propstat")
	err = xw.Start(davNS, "prop")

	err = xw.Start(davNS, "displayname")
	err = xw.Text(false, resource.DisplayName)
	err = xw.End(false)

	if !resource.IsCollection {
		err = xw.Start(davNS, "getcontentlength")
		err = xw.Text(false, strconv.FormatInt(resource.ContentLength, 10))
		err = xw.End(false)
		err = xw.Start(davNS, "getcontenttype")
		err = xw.Text(false, resource.ContentType)
		err = xw.End(false)
	}

	err = xw.Start(davNS, "resourcetype")
	if resource.IsCollection {
		err = xw.Start(davNS, "collection")
		err = xw.End(true)
		err = xw.End(false)
	} else {
		err = xw.End(true)
	}

	err = xw.Start(davNS, "getlastmodified")
	err = xw.Text(false, resource.LastModified)
	err = xw.End(false)
	err = xw.End(false)

	err = xw.Start(davNS, "status")
	err = xw.Text(false, resource.Status)
	err = xw.End(false)

	err = xw.End(false)
	err = xw.End(false)
	return err
}

func MultistatusWebDavResponse(w http.ResponseWriter, response webdav.MultistatusSchema, fromRoot bool) {
	logger := logging.Get("xmlResponse")
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.WriteHeader(207)

	dPrefix := xmlwriter.Prefix("D")
	davNS := xmlwriter.NS("DAV:")

	xw := xmlwriter.New(w)
	xw.Indent(" ")
	err := xw.DefaultProcInst()
	err = xw.Start(davNS, "multistatus", davNS.Bind(dPrefix))

	if fromRoot {
		err = rootCollectionObjectWriter(xw, &davNS)
	}
	for _, resource := range response.Responses {
		err = resourceObjectWriter(xw, &davNS, resource)
	}

	err = xw.End(false)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, "Failed to write XML response", http.StatusInternalServerError)
		return
	}

	_ = xw.Close()
}
