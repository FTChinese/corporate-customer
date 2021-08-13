package main

import (
	"errors"
	"github.com/FTChinese/ftacademy/web"
	rice "github.com/GeertJohan/go.rice"
	"github.com/flosch/pongo2"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"io"
)

// NewRenderer creates a new instance of Renderer based on runtime configuration.
// Deprecated
func NewRenderer(conf web.Config) (Renderer, error) {
	// In debug mode, we use pongo's default local file system loader.
	if conf.Debug {
		log.Info("Development environment using local file system loader")

		loader := pongo2.MustNewLocalFileSystemLoader("web/template")
		set := pongo2.NewSet("local", loader)
		set.Debug = true

		return Renderer{
			config:      conf,
			templateSet: set,
		}, nil
	}

	log.Info("Production environment using rice template loader")
	box, err := rice.FindBox("web/template")
	if err != nil {
		return Renderer{}, err
	}
	loader := NewRiceTemplateLoader(box)
	set := pongo2.NewSet("rice", loader)
	set.Debug = false

	return Renderer{
		config:      conf,
		templateSet: set,
	}, nil
}

// MustNewRenderer panics.
// Deprecated
func MustNewRenderer(config web.Config) Renderer {
	r, err := NewRenderer(config)
	if err != nil {
		log.Fatal(err)
	}

	return r
}

// Renderer is used to render pong2 templates.
// Deprecated.
type Renderer struct {
	templateSet *pongo2.TemplateSet // Load templates from filesystem or rice.
	config      web.Config
}

// Render implements pongo render interface.
// Deprecated.
func (r Renderer) Render(w io.Writer, name string, data interface{}, e echo.Context) error {
	var ctx = pongo2.Context{}

	if data != nil {
		var ok bool
		ctx, ok = data.(pongo2.Context)
		if !ok {
			return errors.New("no pongo2.Context data was passed")
		}
	}

	var t *pongo2.Template
	var err error

	if r.config.Debug {
		// In development the file is loaded from local
		// file system.
		t, err = r.templateSet.FromFile(name)
	} else {
		// In production the file is loaded from rice.
		t, err = r.templateSet.FromCache(name)
	}

	if err != nil {
		return err
	}

	ctx["env"] = r.config

	return t.ExecuteWriter(ctx, w)
}

// RiceTemplateLoader implements pongo2.TemplateLoader to
// loads templates from compiled binary
// Deprecated.
type RiceTemplateLoader struct {
	box *rice.Box
}

// NewRiceTemplateLoader creates a new instance of RiceTemplateLoader.
// Deprecated.
func NewRiceTemplateLoader(box *rice.Box) *RiceTemplateLoader {

	return &RiceTemplateLoader{box: box}
}

func (loader RiceTemplateLoader) Abs(base, name string) string {
	return name
}

func (loader RiceTemplateLoader) Get(path string) (io.Reader, error) {
	return loader.box.Open(path)
}
