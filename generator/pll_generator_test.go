package generator

import (
	"bytes"
	"os"
	"strings"
	"testing"
	"text/template"

	"github.com/komarovn654/embedded_configurator/utils"
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
	utils.InitializeLogger()

	tmpDir, err := os.MkdirTemp("./", "TestCopy")
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
	if err != nil {
		require.Fail(t, "setup environment error")
	}
	defer os.RemoveAll(p.dir)

	t.Run("new generator", func(t *testing.T) {
		gnrt, err := New(p.tmpl, p.dst)
		require.NoError(t, err)
		require.Equal(t, p.dst, gnrt.dstPath)
		assertPllTemplate(t, gnrt, "<no value>")
	})
}

func TestInit(t *testing.T) {
	p, err := setupTmp()
	if err != nil {
		require.Fail(t, "setup environment error")
	}
	defer os.RemoveAll(p.dir)

	t.Run("new template", func(t *testing.T) {
		gnrt := PllGenerator{dstPath: p.dst}
		err = gnrt.init(p.tmpl)
		require.NoError(t, err)
		assertPllTemplate(t, &gnrt, "<no value>")
	})

	t.Run("template does not exist", func(t *testing.T) {
		gnrt := PllGenerator{dstPath: p.dst}
		err = gnrt.init("file does not exist")
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
		gnrt := PllGenerator{dstPath: p.dst, tmpl: tmplt}

		err = gnrt.GenerateHeader("rep string")
		require.NoError(t, err)
		cmpTmpltWithDst(t, p.tmpl, p.dst, "rep string")
	})

	t.Run("generate to stdout", func(t *testing.T) {
		tmplt, err := template.ParseFiles(p.tmpl)
		require.NoError(t, err)
		gnrt := PllGenerator{dstPath: "", tmpl: tmplt}

		err = gnrt.GenerateHeader("rep string")
		require.NoError(t, err)
	})
}
