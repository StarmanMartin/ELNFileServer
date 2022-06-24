package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"os"
	"path"
	"regexp"
	// "fmt"
	// "io"
	"net/http"
)

func checkErr(err error) bool {
	if err != nil {
		ErrorLogger.Println(err)
		return false
	}

	return true
}

type Config struct {
	Root_dir   string `yaml:"root_dir"`
	Prefix_url string `yaml:"webdav_prefix_url"`
	Logfile    string `yaml:"logfile"`
	Host       string `yaml:"host"`
	Admin      string `yaml:"admin_password"`
	Port       int    `yaml:"port"`
}

var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
	cfg         Config
	adminPath   *regexp.Regexp
	webdavPath  *regexp.Regexp
)

// init initializes the logger and parses CMD args.
func init() {
	f, err := os.Open(path.Join("config.yml"))
	if err != nil {
		panic(err)
	}
	defer func(f *os.File) {
		checkErr(f.Close())
	}(f)

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		panic(err)
	}

	re := regexp.MustCompile("^/|/$")
	cfg.Prefix_url = "/" + re.ReplaceAllString(cfg.Prefix_url, "")
	cfg.Host = re.ReplaceAllString(cfg.Host, "")

	logFile, err := os.OpenFile(cfg.Logfile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	mw := io.MultiWriter(os.Stdout, logFile)

	InfoLogger = log.New(mw, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(mw, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	initDb()
	initWebDav()

	adminPath = regexp.MustCompile(fmt.Sprintf(`^%s/?$`, cfg.Prefix_url))
	webdavPath = regexp.MustCompile(fmt.Sprintf(`^%s/?.`, cfg.Prefix_url))

}

func main() {
	http.HandleFunc("/", basicAuth(route))
	InfoLogger.Printf("Server started at %s%s", cfg.Host, cfg.Prefix_url)
	err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), nil)
	checkErr(err)
}

func route(w http.ResponseWriter, r *http.Request, user User) {
	switch {
	case adminPath.MatchString(r.URL.Path):
		AddUserProject(w, r, user)
	case webdavPath.MatchString(r.URL.Path):
		webdavHandler(w, r, user)
	default:
		_, err := w.Write([]byte("Unknown Pattern"))
		checkErr(err)
	}
}
