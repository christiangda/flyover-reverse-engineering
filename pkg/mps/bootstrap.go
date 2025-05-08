package mps

import (
	"os"
	"path"
	"time"

	"github.com/christiangda/flyover-reverse-engineering/pkg/mps/auth"
	"github.com/christiangda/flyover-reverse-engineering/pkg/mps/config"
	"github.com/christiangda/flyover-reverse-engineering/pkg/web"
	"google.golang.org/protobuf/proto"
)

type Context struct {
	AuthContext
	ResourceManifest
}

// Init creates a context from cache or net
func Init(cache Cache, config config.Config) (ctx Context, err error) {
	rm, err := getResourceManifest(cache, config)
	if err != nil {
		return
	}
	s, err := getSession(cache)
	if err != nil {
		return
	}
	authCtx := AuthContext{s, rm, TokenP1(config.TokenP1)}
	ctx = Context{authCtx, rm}
	return
}

// getSession creates a new session or fetches it from cache
func getSession(cache Cache) (s Session, err error) {
	rawSidCachePath := path.Join(cache.Directory, "session.txt")
	var rawSid []byte
	var info os.FileInfo
	if cache.Enabled {
		info, err = os.Stat(rawSidCachePath)
	}
	if !cache.Enabled || os.IsNotExist(err) || err == nil && time.Since(info.ModTime()).Hours() > 24 {
		// from generator
		rawSid = []byte(auth.GenRandStr(40, "0123456789"))
		if cache.Enabled {
			// to cache
			if err = os.WriteFile(rawSidCachePath, rawSid, 0o644); err != nil {
				return
			}
		}
	} else if err == nil {
		// from cache
		if rawSid, err = os.ReadFile(rawSidCachePath); err != nil {
			return
		}
	} else {
		return
	}
	s = Session{ID: string(rawSid)}
	return
}

// getResourceManifest fetches resource manifest from cache or web and decodes it
func getResourceManifest(cache Cache, config config.Config) (rm ResourceManifest, err error) {
	rawRmCachePath := path.Join(cache.Directory, "ResourceManifest.pbd")
	var rawRm []byte
	var info os.FileInfo
	if cache.Enabled {
		info, err = os.Stat(rawRmCachePath)
	}
	if !cache.Enabled || os.IsNotExist(err) || err == nil && time.Since(info.ModTime()).Hours() > 1 {
		// from url
		if rawRm, err = web.Get(config.ResourceManifestURL); err != nil {
			return
		}
		if cache.Enabled {
			// to cache
			if err = os.WriteFile(rawRmCachePath, rawRm, 0o644); err != nil {
				return
			}
		}
	} else if err == nil {
		// from cache
		if rawRm, err = os.ReadFile(rawRmCachePath); err != nil {
			return
		}
	} else {
		return
	}

	// decode resource manifest
	rm = ResourceManifest{}
	if err = proto.Unmarshal(rawRm, &rm); err != nil {
		return
	}

	return
}
