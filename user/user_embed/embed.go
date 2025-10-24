package user_embed

import (
	"embed"

	"server/filefunc"
	"server/srv_conf"
)

//go:embed css/*
//go:embed html/*
//go:embed js/*
//go:embed media/*
var static embed.FS

func User_EmbedFiles() error {

	// write user_embed to disk
	err := filefunc.WriteWebFSToDisk(srv_conf.StaticDir, static)
	if err != nil {
		return err
	}

	return nil

}
