package main

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/regclient/regclient"
	"github.com/regclient/regclient/types/manifest"
	"github.com/regclient/regclient/types/ref"
)

func GetManifest(name string) (manifest.Manifest, error) {
	rc := regclient.New()
	ref, err := ref.New(name)
	if err != nil {
		log.Printf("got err parsing image name manifest: %s", err)
		return nil, err
	}

	m, err := rc.ManifestGet(context.Background(), ref)
	if err != nil {
		log.Printf("got err getting manifest: %s", err)
		return nil, err
	}

	if !m.IsList() {
		err := fmt.Errorf("provided image name has no manifest list")
		log.Print(err)
		return nil, err
	}
	return m, nil
}

func DoesImageSupportArm64(name string) bool {
	m, err := GetManifest(name)
	if err != nil {
		log.Printf("got err getting manifest: %s", err)
		return false
	}

	platforms, err := manifest.GetPlatformList(m)
	if err != nil {
		log.Printf("got err getting platforms for manifest: %s", err)
		return false
	}

	for _, pl := range platforms {
		if pl.String() == "linux/arm64" {
			return true
		}
	}
	return false
}
