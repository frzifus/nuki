package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/frzifus/nuki"
)

var version string

var (
	conf config
)

func init() {
	c := flag.String("-c", "access.json", "config file")
	ip := flag.String("--ip", "", "ip addr")
	port := flag.String("--port", "", "port")
	token := flag.String("--token", "", "token")
	auth := flag.Bool("--auth", false, "set auth to authenticate")
	flag.Parse()

	b, err := ioutil.ReadFile(*c)
	if err == nil {
		if err := json.Unmarshal(b, &conf); err != nil {
			log.Println(*c, "parsing error")
		}
	} else {
		log.Println(*c, "not Found!")
	}

	if *ip != "" {
		conf.IP = *ip
	}
	if *port != "" {
		conf.Port = *port
	}
	if *token != "" {
		conf.Token = *token
	}
	conf.Auth = *auth
}

type config struct {
	IP    string `json:"ip"`
	Port  string `json:"port"`
	Token string `json:"token"`
	Auth  bool
}

func main() {
	logger := log.New()
	logger.Infof("Version %s", version)
	var n *nuki.Nuki
	if conf.Auth {
		n = nuki.NewNuki(conf.IP, conf.Port)
		log.Println("Start Authentication")
		if err := n.Auth; err != nil {
			log.Fatalln("Authentication failed")
		}
		log.Println("Authentication successful, your token is:", n.Token())
	} else {
		if conf.Token == "" {
			log.Fatalln("missing token")
		}
		n = nuki.NewNukiWithToken(conf.IP, conf.Port, conf.Token)
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
