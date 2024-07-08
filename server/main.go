package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"time"

	"github.com/thatisuday/commando"
	"github.com/xeonn/bitlicense"
	"github.com/xeonn/bitlicense/server/fxserver"
	"gitlab.com/bitify-pub/byutils/timeutils"
)

// server port
var port = 9900

var dbUri = "http://localhost:5984"
var dbUser = "admin"
var dbPass = "admin"

//go:embed frontend
var Content embed.FS

func main() {

	now := timeutils.RoundDownTo5Minutes(time.Now().UTC())

	commando.
		SetExecutableName("bitlicense").
		SetVersion("1.0.0").
		SetDescription("Server for issuing license file")

	commando.
		Register(nil).
		SetAction(func(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
			fmt.Println("missing subcommand. see --help for more info")
		})

	commando.
		Register("server").
		SetDescription("Start server that allows license generation and validation via rest api").
		SetShortDescription("start rest api server").
		AddFlag("port,P", "port number for server", commando.Int, port).
		AddFlag("dburi,d", "uri to license database (not yet implemented)", commando.String, dbUri).
		AddFlag("dbuser,u", "database user name (not yet implemented)", commando.String, dbUser).
		AddFlag("dbpass,p", "database user password (not yet implemented)", commando.String, dbPass).
		SetAction(func(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
			// port, _ := flags["port"].GetInt()
			dbUri, _ := flags["dburi"].GetString()
			// dbUser, _ := flags["dbuser"].GetString()
			// dbPass, _ := flags["dbpass"].GetString()

			fmt.Println("Using database path at: " + dbUri)
			startServer()
		})

	commando.
		Register("generate").
		SetDescription(
`
Generate a license based on installed certificate.

Eg: bitlicense generate --expiry 2023-12-06T08:00:00Z --client client1

Expiry date can be in format 2023-12-06T08:00:00Z
or 2023-12-06 (which will be assumed to be 00:00:00)
`,
			).
		SetShortDescription("generate license").
		AddFlag("dburi,d", "uri to license database (not yet implemented)", commando.String, dbUri).
		AddFlag("dbuser,u", "database user name (not yet implemented)", commando.String, dbUser).
		AddFlag("dbpass,p", "database user password (not yet implemented)", commando.String, dbPass).
		AddFlag("client,c", "client name", commando.String, nil).
		AddFlag("expiry,e", "expiry date and time in UTC: format 2023-12-06 or 2023-12-06T08:00:00Z", commando.String, now.Format(time.RFC3339)).
		SetAction(func(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
			// dbUri, _ := flags["dburi"].GetString()
			// dbUser, _ := flags["dbuser"].GetString()
			// dbPass, _ := flags["dbpass"].GetString()
			client, _ := flags["client"].GetString()
			expiry, _ := flags["expiry"].GetString()

			// parse expiry time. Expiry date can be in format 2023-12-06T08:00:00Z
			// or 2023-12-06 (which will be assumed to be 00:00:00)

			expiryTime, err := time.Parse(time.RFC3339, expiry)
			if err != nil {
				expiryTime, err = time.Parse("2006-01-02", expiry)
				if err != nil {
					fmt.Println("Invalid expiry time format")
					return
				}
			}

			// expiry time must be rounded to the previous 5 minutes
			rounded := timeutils.RoundDownTo5Minutes(expiryTime.UTC())

			lic := bitlicense.Issue(client, rounded)

			jstr, err := json.Marshal(lic)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(string(jstr))
		})

	commando.Register("validate").
		SetDescription("Validate a license file").
		SetShortDescription("validate license").
		AddFlag("file,f", "license file", commando.String, nil).
		SetAction(func(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
			file, _ := flags["file"].GetString()
			if file == "" {
				panic("missing file")
			}

			if ok := bitlicense.ValidateFile(file); !ok {
				fmt.Println("License file is invalid")
				return
			} else {
				fmt.Println("License file is valid")
			}
		})

	commando.Parse(nil)

}

func startServer() {
	fxserver.Content = Content
	fmt.Println("Content:", Content)
	fxserver.StartFx(port)
}
