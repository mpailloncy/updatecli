package docker

import (
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/updatecli/updatecli/pkg/plugins/docker/registry/dockerhub"
	"github.com/updatecli/updatecli/pkg/plugins/docker/registry/dockerregistry"
	"github.com/updatecli/updatecli/pkg/plugins/docker/registry/quay"
)

// Source retrieve docker image tag digest from a registry
func (d *Docker) Source(workingDir string) (string, error) {
	errs := d.Validate()

	if len(errs) > 0 {
		for _, err := range errs {
			logrus.Errorln(err)
		}
		return "", errors.New("error found in docker parameters")
	}

	hostname, image, err := parseImage(d.Image)

	if err != nil {
		return "", err
	}

	var r Registry

	if d.isDockerHub() {
		dh := dockerhub.Docker{
			Image: image,
			Tag:   d.Tag,
			Token: d.Token,
		}

		r = &dh

	} else if d.isQuaiIO() {

		q := quay.Docker{
			Image: image,
			Tag:   d.Tag,
			Token: d.Token,
		}

		r = &q

	} else if ok, err := d.IsDockerRegistry(); ok {
		if err != nil {
			return "", err
		}

		dr := dockerregistry.Docker{
			Image:    image,
			Tag:      d.Tag,
			Hostname: hostname,
			Token:    d.Token,
		}

		r = &dr

	} else {
		return "", fmt.Errorf("unknown docker registry api")
	}

	digest, err := r.Digest()

	if err != nil {
		return "", err
	}

	if digest == "" {
		logrus.Infof("\u2717 No Digest found for docker image %s:%s on the Docker Registry", d.Image, d.Tag)
	} else {
		logrus.Infof("\u2714 Digest '%v' found for docker image %s:%s available from Docker Registry", digest, d.Image, d.Tag)
		logrus.Infof("Remark: Do not forget to add @sha256 after your the docker image name")
		logrus.Infof("Example: %v@sha256:%v", d.Image, digest)
	}

	return digest, nil
}
