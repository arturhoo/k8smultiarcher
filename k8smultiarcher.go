package main

import (
	"github.com/arturhoo/k8smultiarcher/image"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	name := "gcr.io/distroless/static-debian12"
	arm64Supported := image.DoesImageSupportArm64(name)
	if arm64Supported {
		log.Info().Msg("arm is supported!")
	} else {
		log.Info().Msg("arm not supported!")
	}
}
