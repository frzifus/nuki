package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/frzifus/nuki"
)

var version string

type config struct {
	Address string `json:"address"`
	Token   string `json:"token"`
	Auth    bool
}

type raw struct{ logger log.FieldLogger }

func (r raw) Do(req *http.Request) (*http.Response, error) {
	b, err := httputil.DumpRequest(req, true)
	if err != nil {
		return nil, err
	}
	r.logger.Infof("raw:\n%s", string(b))
	os.Exit(0)
	return &http.Response{}, nil
}

func main() {
	var conf config
	flag.StringVar(&conf.Address, "address", "", "address e.g.: 192.168.1.100:1234")
	flag.StringVar(&conf.Token, "token", "", "token e.g.: 12345")
	flag.BoolVar(&conf.Auth, "auth", false, "set auth to get a token")
	c := flag.String("-c", "", "config file")
	levelDebug := flag.Bool("debug", false, "use debug log level")
	dump := flag.Bool("dump", false, "print http request")
	flag.Parse()

	logger := log.New()
	if *levelDebug {
		logger.SetLevel(log.DebugLevel)
	}

	if *c != "" {
		b, err := ioutil.ReadFile(*c)
		if err != nil {
			logger.Fatalln(err)
		}
		if err := json.Unmarshal(b, &conf); err != nil {
			logger.Fatalln(*c, "parsing error")
		}
	}

	logger.Infof("Version %s", version)
	var opts []nuki.Option
	if conf.Token != "" {
		opts = append(opts, nuki.WithToken(conf.Token))
	}
	if *dump {
		opts = append(opts, nuki.WithHTTPClient(&raw{logger: logger}))
	}
	n := nuki.NewNuki(conf.Address, opts...)
	if conf.Auth {
		logger.Debug("Start Authentication")
		if _, err := n.Auth(); err != nil {
			logger.Fatalf("Authentication failed: %s", err.Error())
		}
		logger.Info("Authentication successful, your token is:", n.Token())
		return
	}

	commands := func() {
		fmt.Println("cmd:")
		fmt.Println("  list")
		fmt.Println("  unlatch <nukiID>")
		fmt.Println("  lock    <nukiID>")
		fmt.Println("  unlock  <nukiID>")
	}

	args := flag.Args()
	if len(args) < 1 {
		commands()
		os.Exit(0)
	}

	switch args[0] {
	case "list":
		list, err := n.List()
		if err != nil {
			logger.Fatalln(err)
		}
		for _, x := range list {
			logger.WithFields(log.Fields{
				"action":          "list",
				"nukiID":          x.NukiID,
				"name":            x.Name,
				"state":           x.LastKnownState.StateName,
				"batteryCritical": x.LastKnownState.BatteryCritical,
			}).Info()
		}

	case "unlatch":
		if len(args) < 2 {
			logger.Println("./nuki-cli", args[0], "<nukiID>")
			os.Exit(0)
		}
		nukiID, err := strconv.Atoi(args[1])
		if err != nil {
			logger.Fatal(err)
		}
		res, err := n.LockAction(nukiID, nuki.ActionUnlatch, false)
		if err != nil {
			logger.Fatal(err)
		}
		logger.WithFields(log.Fields{
			"action":          "unlatch",
			"nukiID":          nukiID,
			"success":         res.Success,
			"batteryCritical": res.BatteryCritical,
		}).Info()

	case "lock":
		if len(args) < 2 {
			logger.Println("./nuki-cli", args[0], "<nukiID>")
			os.Exit(0)
		}

		nukiID, err := strconv.Atoi(flag.Args()[1])
		if err != nil {
			logger.Fatal(err)
		}
		res, err := n.LockAction(nukiID, nuki.ActionLock, false)
		if err != nil {
			logger.Fatal(err)
		}
		logger.WithFields(log.Fields{
			"action":          "lock",
			"nukiID":          nukiID,
			"success":         res.Success,
			"batteryCritical": res.BatteryCritical,
		}).Info()

	case "unlock":
		if len(args) < 2 {
			fmt.Println("./nuki-cli", args[0], "<nukiID>")
			os.Exit(0)
		}
		nukiID, err := strconv.Atoi(flag.Args()[1])
		if err != nil {
			log.Fatal(err)
		}
		res, err := n.LockAction(nukiID, nuki.ActionUnlock, false)
		if err != nil {
			log.Fatal(err)
		}

		logger.WithFields(log.Fields{
			"action":          "unlock",
			"nukiID":          nukiID,
			"success":         res.Success,
			"batteryCritical": res.BatteryCritical,
		}).Info()

	default:
		fmt.Println("unknown command")
		commands()
	}

}
