package jumpsuit

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"net/http"
	"net/http/httputil"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/gertd/go-pluralize"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

var (
	plc = pluralize.NewClient()
)

type Server struct {
	EchoServer *echo.Echo
	Storage    Storage
}

type ServerOpts struct {
	Storage Storage
}

func New(opts *ServerOpts) *Server {
	es := echo.New()
	es.Logger.SetLevel(log.DEBUG)
	return &Server{
		EchoServer: es,
		Storage:    opts.Storage,
	}
}

type APIOptions struct {
	Path string
}

func (s *Server) NewAPI(t interface{}, opts APIOptions) {
	objType := reflect.TypeOf(t)
	if opts.Path == "" {
		opts.Path = strings.ToLower(objType.Name())
	}
	if plc.IsSingular(opts.Path) {
		opts.Path = plc.Pluralize(opts.Path, 2, false)
	}
	s.EchoServer.Any("api/"+opts.Path, s.basePath(opts.Path))
	s.EchoServer.Any("api/"+opts.Path+"/", s.basePath(opts.Path))
	s.EchoServer.Any("api/"+opts.Path+"/:id", s.selectedPath(opts.Path))
}

func (s *Server) basePath(table string) echo.HandlerFunc {
	return func(c echo.Context) error {
		// log request data
		b, err := httputil.DumpRequest(c.Request(), true)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(b))
		switch c.Request().Method {
		case http.MethodGet:
			fmt.Println("getting")
			return s.Lst(table, c)
		case http.MethodPut:
			return s.Put(table, c)
		}
		return nil
	}
}

func (s *Server) Lst(table string, c echo.Context) error {
	lst, err := s.Storage.Lst(table)
	if err != nil {
		fmt.Println("1", err)
		return err
	}
	raw, err := json.Marshal(lst)
	if err != nil {
		fmt.Println("2", err)
		return err
	}
	return c.Blob(200, echo.MIMEApplicationJSON, raw)
}

func (s *Server) selectedPath(table string) echo.HandlerFunc {
	return func(c echo.Context) error {
		// log request data
		b, err := httputil.DumpRequest(c.Request(), true)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(b))

		switch c.Request().Method {
		case http.MethodGet:
			return s.Get(table, c)
		case http.MethodDelete:
			return s.Del(table, c)
		case http.MethodPut:
			return s.Put(table, c)
		}
		return nil
	}
}

func (s *Server) Get(table string, c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return err
	}
	obj, err := s.Storage.Get(table, id)
	if err != nil {
		return err
	}
	raw, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	return c.Blob(200, echo.MIMEApplicationJSON, raw)
}

func (s *Server) Del(table string, c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return err
	}
	err = s.Storage.Del(table, id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusAccepted, nil)
}

func (s *Server) Put(table string, c echo.Context) error {
	obj := new(map[string]any)
	b, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, obj)
	if err != nil {
		return err
	}

	var (
		idParam = c.Param("id")
		id      int64
	)

	if idParam == "" {
		id, err = s.Storage.Inc(table)
		if err != nil {
			return err
		}
		(*obj)["ID"] = id
	} else {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			return err
		}
		objID := (*obj)["ID"].(float64)
		if int64(objID) != id {
			return fmt.Errorf("invalid id \"ID\"")
		}
		if _, err := s.Storage.Get(table, id); err != nil {
			return err
		}
	}

	err = s.Storage.Put(table, id, obj)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusAccepted, map[string]int64{"updated": id})
}

func (s *Server) Start(addr string) error {
	return s.EchoServer.Start(addr)
}

func fields(t interface{}) []reflect.StructField {
	objType := reflect.TypeOf(t)
	l := objType.NumField()
	ret := []reflect.StructField{}
	for i := 0; i < l; i++ {
		ret = append(ret, objType.Field(i))
	}
	return ret
}

type obj struct {
	Name   string
	Fields []reflect.StructField
}

func Selected(t interface{}, wr io.Writer) error {
	tmpl := template.Must(template.ParseFiles("tmpl/selected.go.html"))
	objType := reflect.TypeOf(t)
	tmplObj := obj{
		Name:   objType.Name(),
		Fields: fields(t),
	}

	return tmpl.Execute(wr, tmplObj)
}

func NewFileStore(path string) *FileStore {
	_, err := os.Stat(path)
	if err != nil {
		os.MkdirAll(path, fs.ModePerm)
	}

	return &FileStore{
		Mutex:    sync.Mutex{},
		BasePath: path,
	}
}

func NewMemoryStore() *Memory {
	return &Memory{
		Mutex: sync.Mutex{},
		Data:  make(map[int64]any),
	}
}
