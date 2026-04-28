// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package cmd

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/encoding/gyaml"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/os/gview"
	"github.com/gogf/gf/v2/util/gvalid"

	"github.com/gogf/gf/cmd/gf/v2/internal/utility/mlog"
)

//go:embed resources/i18n/*.yaml
var i18nFS embed.FS

//go:embed resources/templates/index.html
var configEditorHTML string

//go:embed resources/static/*
var staticFS embed.FS

var (
	// CfgEditor is the management object for `gf config` command.
	CfgEditor = cCfgEditor{}
)

type cCfgEditor struct {
	g.Meta `name:"config" brief:"start the configuration visual editor"`
}

type cCfgEditorInput struct {
	g.Meta `name:"config" config:"gfcli.config"`
	Port   int    `short:"p" name:"port"   brief:"web server port" d:"8888"`
	File   string `short:"f" name:"file"   brief:"configuration file path"`
}

type cCfgEditorOutput struct{}

func init() {
	registerAllSchemas()
}

// registerAllSchemas registers configuration schemas for the five core modules.
func registerAllSchemas() {
	// Server
	gcfg.RegisterSchema("server", "server", ghttp.ServerConfig{}, map[string]string{
		"Name": "Basic", "Address": "Basic", "HTTPSAddr": "Basic",
		"HTTPSCertPath": "Basic", "HTTPSKeyPath": "Basic",
		"ReadTimeout": "Basic", "WriteTimeout": "Basic", "IdleTimeout": "Basic",
		"MaxHeaderBytes": "Basic", "KeepAlive": "Basic", "ServerAgent": "Basic",
		"IndexFolder": "Static", "ServerRoot": "Static", "FileServerEnabled": "Static",
		"CookieMaxAge": "Cookie", "CookiePath": "Cookie", "CookieDomain": "Cookie",
		"CookieSameSite": "Cookie", "CookieSecure": "Cookie", "CookieHttpOnly": "Cookie",
		"SessionIdName": "Session", "SessionMaxAge": "Session", "SessionPath": "Session",
		"SessionCookieMaxAge": "Session", "SessionCookieOutput": "Session",
		"LogPath": "Logging", "LogLevel": "Logging", "LogStdout": "Logging",
		"ErrorStack": "Logging", "ErrorLogEnabled": "Logging", "ErrorLogPattern": "Logging",
		"AccessLogEnabled": "Logging", "AccessLogPattern": "Logging",
		"PProfEnabled": "PProf", "PProfPattern": "PProf",
		"OpenApiPath": "API", "SwaggerPath": "API", "SwaggerUITemplate": "API",
		"Graceful": "Graceful", "GracefulTimeout": "Graceful", "GracefulShutdownTimeout": "Graceful",
		"ClientMaxBodySize": "Other", "FormParsingMemory": "Other",
		"NameToUriType": "Other", "RouteOverWrite": "Other", "DumpRouterMap": "Other",
		"Endpoints": "Other", "Rewrites": "Other", "IndexFiles": "Other", "SearchPaths": "Other",
		"StaticPaths": "Other", "Listeners": "Other",
	})

	// Database
	gcfg.RegisterSchema("database", "database", gdb.ConfigNode{}, map[string]string{
		"Host": "Connection", "Port": "Connection", "User": "Connection",
		"Pass": "Connection", "Name": "Connection", "Type": "Connection",
		"Link": "Connection", "Extra": "Connection", "Protocol": "Connection",
		"Charset": "Connection", "Timezone": "Connection", "Namespace": "Connection",
		"MaxIdleConnCount": "Pool", "MaxOpenConnCount": "Pool",
		"MaxConnLifeTime": "Pool", "MaxIdleConnTime": "Pool",
		"Role": "Role", "Debug": "Role", "Prefix": "Role", "DryRun": "Role", "Weight": "Role",
		"QueryTimeout": "Timeout", "ExecTimeout": "Timeout",
		"TranTimeout": "Timeout", "PrepareTimeout": "Timeout",
		"CreatedAt": "AutoTimestamp", "UpdatedAt": "AutoTimestamp",
		"DeletedAt": "AutoTimestamp", "TimeMaintainDisabled": "AutoTimestamp",
	})

	// Redis
	gcfg.RegisterSchema("redis", "redis", gredis.Config{}, map[string]string{
		"Address": "Connection", "Db": "Connection", "User": "Connection",
		"Pass": "Connection", "Protocol": "Connection",
		"MinIdle": "Pool", "MaxIdle": "Pool", "MaxActive": "Pool",
		"MaxConnLifetime": "Pool", "IdleTimeout": "Pool", "WaitTimeout": "Pool",
		"DialTimeout": "Timeout", "ReadTimeout": "Timeout", "WriteTimeout": "Timeout",
		"MasterName": "Sentinel", "SentinelUser": "Sentinel", "SentinelPass": "Sentinel",
		"TLS": "Security", "TLSSkipVerify": "Security",
		"SlaveOnly": "Security", "Cluster": "Security",
	})

	// Logger
	gcfg.RegisterSchema("logger", "logger", glog.Config{}, map[string]string{
		"Flags": "Basic", "TimeFormat": "Basic", "Path": "Basic",
		"File": "Basic", "Level": "Basic", "Prefix": "Basic",
		"HeaderPrint": "Output", "StdoutPrint": "Output", "LevelPrint": "Output",
		"StdoutColorDisabled": "Output", "WriterColorEnable": "Output",
		"StSkip": "Stack", "StStatus": "Stack", "StFilter": "Stack",
		"RotateSize": "Rotate", "RotateExpire": "Rotate",
		"RotateBackupLimit": "Rotate", "RotateBackupExpire": "Rotate",
		"RotateBackupCompress": "Rotate", "RotateCheckInterval": "Rotate",
	})

	// Viewer
	gcfg.RegisterSchema("viewer", "viewer", gview.Config{}, map[string]string{
		"Paths": "Basic", "Data": "Basic", "DefaultFile": "Basic",
		"Delimiters": "Basic", "AutoEncode": "Basic",
	})
}

// Index starts the config editor web server.
func (c cCfgEditor) Index(ctx context.Context, in cCfgEditorInput) (out *cCfgEditorOutput, err error) {
	mlog.Printf("[ConfigEditor] Starting with port=%d, file=%q", in.Port, in.File)

	// Verify embedded i18n files are accessible.
	for _, lang := range []string{"en", "zh-CN"} {
		path := "resources/i18n/" + lang + ".yaml"
		if data, e := i18nFS.ReadFile(path); e != nil {
			mlog.Printf("[ConfigEditor] WARNING: embedded i18n file %q not found: %v", path, e)
		} else {
			mlog.Printf("[ConfigEditor] Embedded i18n file %q loaded, size=%d bytes", path, len(data))
		}
	}

	s := g.Server("gf-config-editor")
	s.SetPort(in.Port)
	s.SetDumpRouterMap(false)

	// API endpoints.
	s.Group("/api", func(group *ghttp.RouterGroup) {
		group.GET("/schemas", apiGetSchemas)
		group.GET("/config", apiGetConfig(in.File))
		group.POST("/config/validate", apiValidateConfig)
		group.POST("/config/save", apiSaveConfig)
		group.GET("/i18n/:lang", apiGetI18n)
	})

	// Serve embedded static files.
	s.BindHandler("/static/*", func(r *ghttp.Request) {
		filePath := strings.TrimPrefix(r.URL.Path, "/static/")
		data, err := fs.ReadFile(staticFS, "resources/static/"+filePath)
		if err != nil {
			r.Response.WriteStatus(http.StatusNotFound)
			return
		}
		if strings.HasSuffix(filePath, ".js") {
			r.Response.Header().Set("Content-Type", "application/javascript; charset=utf-8")
		} else if strings.HasSuffix(filePath, ".css") {
			r.Response.Header().Set("Content-Type", "text/css; charset=utf-8")
		}
		r.Response.Header().Set("Cache-Control", "public, max-age=86400")
		r.Response.Write(data)
	})

	// Serve the embedded UI.
	s.BindHandler("/", func(r *ghttp.Request) {
		r.Response.WriteHeader(http.StatusOK)
		r.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
		r.Response.Write(configEditorHTML)
	})

	addr := fmt.Sprintf("http://127.0.0.1:%d", in.Port)
	mlog.Printf("[ConfigEditor] GoFrame Config Editor starting at %s", addr)

	go func() {
		time.Sleep(500 * time.Millisecond)
		if err := openBrowser(addr); err != nil {
			mlog.Printf("[ConfigEditor] WARNING: failed to open browser: %v", err)
		}
	}()

	s.Run()
	return
}

// apiGetSchemas returns all registered module schemas.
func apiGetSchemas(r *ghttp.Request) {
	schemas := gcfg.GetAllSchemas()
	r.Response.WriteJsonExit(g.Map{
		"code": 0,
		"data": schemas,
	})
}

// apiGetConfig returns the current configuration values.
func apiGetConfig(file string) func(r *ghttp.Request) {
	return func(r *ghttp.Request) {
		configFile := file
		if configFile == "" {
			searchPaths := []string{
				"config.yaml", "config.yml", "config.toml", "config.json",
				"config/config.yaml", "config/config.yml",
				"config/config.toml", "config/config.json",
				"manifest/config/config.yaml", "manifest/config/config.yml",
				"manifest/config/config.toml", "manifest/config/config.json",
				"app.yaml", "app.yml",
			}
			for _, name := range searchPaths {
				if gfile.Exists(name) {
					configFile = name
					break
				}
			}
		}

		data := g.Map{}
		filePath := ""
		fileType := ""
		if configFile != "" && gfile.Exists(configFile) {
			filePath = gfile.RealPath(configFile)
			fileType = gfile.ExtName(configFile)
			content := gfile.GetBytes(configFile)
			j, err := gjson.LoadContent(content)
			if err != nil {
				r.Response.WriteJsonExit(g.Map{
					"code":    1,
					"message": fmt.Sprintf("Failed to parse config file %q: %v", filePath, err),
				})
				return
			}
			data = j.Map()
		}

		r.Response.WriteJsonExit(g.Map{
			"code": 0,
			"data": g.Map{
				"config":   data,
				"filePath": filePath,
				"fileType": fileType,
			},
		})
	}
}

// apiValidateConfig validates configuration values using gvalid.
func apiValidateConfig(r *ghttp.Request) {
	var reqData struct {
		Module string         `json:"module"`
		Values map[string]any `json:"values"`
	}
	if err := r.Parse(&reqData); err != nil {
		r.Response.WriteJsonExit(g.Map{"code": 1, "message": err.Error()})
		return
	}

	schema, ok := gcfg.GetSchema(reqData.Module)
	if !ok {
		r.Response.WriteJsonExit(g.Map{"code": 1, "message": fmt.Sprintf("module %q not found", reqData.Module)})
		return
	}

	// Build validation rules from schema fields.
	var rules []string
	for _, field := range schema.Fields {
		if field.Rule == "" {
			continue
		}
		rule := field.JsonKey + "|" + field.Rule
		rules = append(rules, rule)
	}

	if len(rules) > 0 {
		if err := gvalid.New().Data(reqData.Values).Rules(rules).Run(r.Context()); err != nil {
			// Parse validation errors into field-level messages.
			validationErrors := make(map[string]string)
			if vErr, ok := err.(gvalid.Error); ok {
				for _, item := range vErr.Items() {
					for field, ruleErrMap := range item {
						for _, ruleErr := range ruleErrMap {
							validationErrors[field] = ruleErr.Error()
							break
						}
					}
				}
			} else {
				validationErrors["_general"] = err.Error()
			}
			r.Response.WriteJsonExit(g.Map{
				"code":    1,
				"message": "Validation failed",
				"errors":  validationErrors,
			})
			return
		}
	}

	r.Response.WriteJsonExit(g.Map{
		"code":    0,
		"message": "Valid",
	})
}

// apiSaveConfig saves configuration to file.
func apiSaveConfig(r *ghttp.Request) {
	var reqData struct {
		Config   map[string]any `json:"config"`
		FilePath string         `json:"filePath"`
		FileType string         `json:"fileType"`
	}
	if err := r.Parse(&reqData); err != nil {
		r.Response.WriteJsonExit(g.Map{"code": 1, "message": err.Error()})
		return
	}

	if reqData.FilePath == "" {
		reqData.FilePath = "config.yaml"
		reqData.FileType = "yaml"
	}

	var err error
	switch reqData.FileType {
	case "yaml", "yml":
		err = saveYAMLPreservingComments(reqData.FilePath, reqData.Config)
	default:
		j := gjson.New(reqData.Config)
		var content string
		switch reqData.FileType {
		case "toml":
			content, err = j.ToTomlString()
		case "json":
			content, err = j.ToJsonIndentString()
		case "ini":
			content, err = j.ToIniString()
		default:
			content, err = j.ToYamlString()
		}
		if err == nil {
			err = gfile.PutContents(reqData.FilePath, content)
		}
	}

	if err != nil {
		r.Response.WriteJsonExit(g.Map{"code": 1, "message": err.Error()})
		return
	}

	r.Response.WriteJsonExit(g.Map{
		"code":    0,
		"message": "Configuration saved successfully",
		"data": g.Map{
			"filePath": gfile.RealPath(reqData.FilePath),
		},
	})
}

// saveYAMLPreservingComments writes the config map to a YAML file while preserving
// any existing comments in the file.
func saveYAMLPreservingComments(filePath string, newConfig map[string]any) error {
	var (
		docNode yaml.Node
		indent  = 2
	)

	if gfile.Exists(filePath) {
		content := gfile.GetBytes(filePath)
		indent = detectYAMLIndent(content)
		if err := yaml.Unmarshal(content, &docNode); err != nil {
			docNode = yaml.Node{}
		}
	}

	if docNode.Kind == 0 {
		docNode = yaml.Node{Kind: yaml.DocumentNode}
		docNode.Content = []*yaml.Node{{Kind: yaml.MappingNode, Tag: "!!map"}}
	} else if docNode.Kind == yaml.DocumentNode {
		if len(docNode.Content) == 0 {
			docNode.Content = []*yaml.Node{{Kind: yaml.MappingNode, Tag: "!!map"}}
		} else if docNode.Content[0].Kind != yaml.MappingNode {
			docNode.Content = []*yaml.Node{{Kind: yaml.MappingNode, Tag: "!!map"}}
		}
	}

	applyMapToYAMLNode(docNode.Content[0], newConfig)

	var buf strings.Builder
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(indent)
	if err := enc.Encode(&docNode); err != nil {
		return err
	}
	_ = enc.Close()
	return gfile.PutContents(filePath, buf.String())
}

// detectYAMLIndent returns the number of spaces used for indentation in the YAML content.
func detectYAMLIndent(content []byte) int {
	for _, line := range strings.Split(string(content), "\n") {
		trimmed := strings.TrimLeft(line, " ")
		if len(trimmed) == 0 || strings.HasPrefix(trimmed, "#") {
			continue
		}
		spaces := len(line) - len(trimmed)
		if spaces > 0 {
			return spaces
		}
	}
	return 2
}

// applyMapToYAMLNode recursively merges updates into an existing yaml.MappingNode,
// preserving comments and formatting style on nodes that already exist.
func applyMapToYAMLNode(mappingNode *yaml.Node, updates map[string]any) {
	if mappingNode.Kind != yaml.MappingNode {
		return
	}
	keyIndex := make(map[string]int)
	for i := 0; i < len(mappingNode.Content)-1; i += 2 {
		keyIndex[mappingNode.Content[i].Value] = i + 1
	}

	for key, value := range updates {
		valIdx, exists := keyIndex[key]
		if !exists {
			keyNode := &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: key}
			valNode := anyToYAMLNode(value)
			mappingNode.Content = append(mappingNode.Content, keyNode, valNode)
			keyIndex[key] = len(mappingNode.Content) - 1
		} else {
			existingVal := mappingNode.Content[valIdx]
			updateYAMLNodeInPlace(existingVal, value)
		}
	}
}

// updateYAMLNodeInPlace updates the yaml.Node in place to reflect newValue
// while maximally preserving the original formatting style and comments.
func updateYAMLNodeInPlace(node *yaml.Node, newValue any) {
	head, line, foot := node.HeadComment, node.LineComment, node.FootComment

	switch v := newValue.(type) {
	case map[string]any:
		if node.Kind == yaml.MappingNode {
			applyMapToYAMLNode(node, v)
			return
		}
		*node = *anyToYAMLNode(v)

	case []any:
		if node.Kind == yaml.SequenceNode {
			style := node.Style
			newSeq := anyToYAMLNode(v)
			*node = *newSeq
			node.Style = style
		} else {
			*node = *anyToYAMLNode(v)
		}

	default:
		newNode := anyToYAMLNode(v)
		if node.Kind == yaml.ScalarNode && newNode.Kind == yaml.ScalarNode {
			node.Value = newNode.Value
			node.Tag = newNode.Tag
		} else {
			*node = *newNode
		}
	}

	node.HeadComment, node.LineComment, node.FootComment = head, line, foot
}

// anyToYAMLNode converts a Go value to a yaml.Node.
func anyToYAMLNode(v any) *yaml.Node {
	if v == nil {
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!null", Value: "null"}
	}
	switch val := v.(type) {
	case map[string]any:
		node := &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}
		for k, vv := range val {
			keyNode := &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: k}
			valNode := anyToYAMLNode(vv)
			node.Content = append(node.Content, keyNode, valNode)
		}
		return node
	case []any:
		node := &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}
		for _, item := range val {
			node.Content = append(node.Content, anyToYAMLNode(item))
		}
		return node
	case bool:
		s := "false"
		if val {
			s = "true"
		}
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!bool", Value: s}
	case int:
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!int", Value: strconv.Itoa(val)}
	case int64:
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!int", Value: strconv.FormatInt(val, 10)}
	case float64:
		s := strconv.FormatFloat(val, 'f', -1, 64)
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!float", Value: s}
	case string:
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: val}
	default:
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: fmt.Sprintf("%v", v)}
	}
}

// apiGetI18n returns i18n translations for the given language.
func apiGetI18n(r *ghttp.Request) {
	lang := r.Get("lang").String()
	if lang == "" {
		lang = "en"
	}

	fileName := lang + ".yaml"
	filePath := "resources/i18n/" + fileName

	content, err := i18nFS.ReadFile(filePath)
	if err != nil {
		r.Response.WriteJsonExit(g.Map{
			"code": 0,
			"data": g.Map{},
		})
		return
	}

	var translations map[string]string
	if err = gyaml.DecodeTo(content, &translations); err != nil {
		r.Response.WriteJsonExit(g.Map{
			"code": 0,
			"data": g.Map{},
		})
		return
	}

	r.Response.WriteJsonExit(g.Map{
		"code": 0,
		"data": translations,
	})
}

// openBrowser opens the default browser to the given URL.
func openBrowser(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	return cmd.Start()
}
