package pll_generator

import (
	"bytes"
	"os"
	"strings"
	"testing"
	"text/template"

	l "github.com/komarovn654/embedded_configurator/utils/log"
	"github.com/stretchr/testify/require"
)

var (
	TestTemplate = "test template {{.}}"
)

type paths struct {
	dir  string
	tmpl string
	dst  string
}

func setupTmp() (paths, error) {
	l.InitializeLogger()

	tmpDir, err := os.MkdirTemp("./", "Test")
	if err != nil {
		return paths{}, err
	}
	tmpl, err := os.CreateTemp(tmpDir, "TestTmpl")
	if err != nil {
		return paths{}, err
	}
	tmpl.WriteString(TestTemplate)
	dst, err := os.CreateTemp(tmpDir, "TestOut")
	if err != nil {
		return paths{}, err
	}
	tmpl.Close()
	dst.Close()

	return paths{dir: tmpDir, tmpl: tmpl.Name(), dst: dst.Name()}, err
}

func assertPllTemplate(t *testing.T, gnrt *PllGenerator, replaced string) {
	buf := bytes.NewBufferString("")
	err := gnrt.tmpl.Execute(buf, nil)
	require.NoError(t, err)

	require.Equal(t, bytes.NewBufferString(strings.Replace(TestTemplate, "{{.}}", replaced, -1)), buf)
}

func cmpTmpltWithDst(t *testing.T, tmpl string, dst string, replaced string) {
	tmplBytes, err := os.ReadFile(tmpl)
	require.NoError(t, err)
	dstBytes, err := os.ReadFile(dst)
	require.NoError(t, err)

	genTmpl := strings.Replace(string(tmplBytes), "{{.}}", replaced, -1)
	require.Equal(t, genTmpl, string(dstBytes))
}

func TestNew(t *testing.T) {
	p, err := setupTmp()
	require.NoError(t, err, "setup environment error")
	defer os.RemoveAll(p.dir)

	t.Run("new generator", func(t *testing.T) {
		gnrt, err := New(p.tmpl, p.dst)
		require.NoError(t, err)
		require.NotNil(t, gnrt)
		require.Equal(t, p.dst, gnrt.destPath)
		require.Equal(t, p.tmpl, gnrt.tmplPath)
		assertPllTemplate(t, gnrt, "<no value>")
	})
}

func TestInit(t *testing.T) {
	p, err := setupTmp()
	require.NoError(t, err, "setup environment error")
	defer os.RemoveAll(p.dir)

	t.Run("new template", func(t *testing.T) {
		gnrt := PllGenerator{destPath: p.dst, tmplPath: p.tmpl}
		err = gnrt.init()
		require.NoError(t, err)
		assertPllTemplate(t, &gnrt, "<no value>")
	})

	t.Run("template does not exist", func(t *testing.T) {
		gnrt := PllGenerator{destPath: p.dst}
		err = gnrt.init()
		require.Error(t, err)
	})
}

func TestGenerateHeader(t *testing.T) {
	p, err := setupTmp()
	if err != nil {
		require.Fail(t, "setup environment error")
	}
	defer os.RemoveAll(p.dir)

	t.Run("generate to file", func(t *testing.T) {
		tmplt, err := template.ParseFiles(p.tmpl)
		require.NoError(t, err)
		gnrt := PllGenerator{destPath: p.dst, tmpl: tmplt}

		err = gnrt.GenerateHeader("rep string")
		require.NoError(t, err)
		cmpTmpltWithDst(t, p.tmpl, p.dst, "rep string")
	})

	t.Run("generate to stdout", func(t *testing.T) {
		tmplt, err := template.ParseFiles(p.tmpl)
		require.NoError(t, err)
		gnrt := PllGenerator{destPath: "", tmpl: tmplt}

		err = gnrt.GenerateHeader("rep string")
		require.NoError(t, err)
	})
}
