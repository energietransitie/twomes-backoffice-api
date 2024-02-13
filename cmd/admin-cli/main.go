package main

import (
	"fmt"
	"net/rpc"
	"os"
	"text/tabwriter"
	"time"

	"github.com/energietransitie/twomes-backoffice-api/twomes/admin"
	"github.com/urfave/cli/v2"
)

func main() {
	nameFlag := &cli.StringFlag{
		Name:     "name",
		Aliases:  []string{"n"},
		Usage:    "`NAME` of the admin",
		Required: true,
	}

	year, month, day := time.Now().UTC().Date()
	defaultExpiry := fmt.Sprintf("%d-%02d-%02d", year+1, month, day)

	expiryFlag := &cli.StringFlag{
		Name:    "expiry",
		Aliases: []string{"e"},
		Usage:   "`EXPIRATION` date (yyyy-mm-dd) (at 00:00 UTC) of the admin",
		Value:   defaultExpiry,
	}

	app := &cli.App{
		Name:  "admin-cli",
		Usage: "Use this CLI tool to manage API admins",
		Commands: []*cli.Command{
			{
				Name:   "list",
				Usage:  "List all admins",
				Action: listAdmins,
			},
			{
				Name:  "create",
				Usage: "Create a new admin",
				Flags: []cli.Flag{
					nameFlag,
					expiryFlag,
				},
				Action: createAdmin,
			},
			{
				Name:  "delete",
				Usage: "Delete an admin",
				Flags: []cli.Flag{
					nameFlag,
				},
				Action: deleteAdmin,
			},
			{
				Name:  "reactivate",
				Usage: "Reactive an admin",
				Flags: []cli.Flag{
					nameFlag,
				},
				Action: reactivateAdmin,
			},
			{
				Name:  "expiry",
				Usage: "Set expiry for an admin",
				Flags: []cli.Flag{
					nameFlag,
					expiryFlag,
				},
				Action: setAdminExpiry,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}

func listAdmins(ctx *cli.Context) error {
	var admins []admin.Admin

	client, err := getRPCClient()
	if err != nil {
		return err
	}

	err = client.Call("AdminHandler.List", 0, &admins)
	if err != nil {
		return err
	}

	w := new(tabwriter.Writer)

	w.Init(os.Stdout, 4, 4, 4, ' ', 0)

	defer w.Flush()

	fmt.Fprintf(w, "ID\tName\tActivated at\tExpires at\n")

	timeFormat := "2006-01-02 15:04:05 MST"
	for _, admin := range admins {
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\t\n", admin.ID, admin.Name, admin.ActivatedAt.Format(timeFormat), admin.Expiry.Format(timeFormat))
	}

	return nil
}

func createAdmin(ctx *cli.Context) error {
	expiryString := ctx.String("expiry")
	var expiry time.Time
	if expiryString != "" {
		var err error
		expiry, err = time.Parse("2006-01-02", expiryString)
		if err != nil {
			return err
		}
	}

	admin := admin.Admin{
		Name:   ctx.String("name"),
		Expiry: expiry,
	}

	client, err := getRPCClient()
	if err != nil {
		return err
	}

	var authToken string
	err = client.Call("AdminHandler.Create", admin, &authToken)
	if err != nil {
		return err
	}

	fmt.Printf("Admin \"%s\" created. Authorization token: %s\n", admin.Name, authToken)
	return nil
}

func deleteAdmin(ctx *cli.Context) error {
	admin := admin.Admin{
		Name: ctx.String("name"),
	}

	client, err := getRPCClient()
	if err != nil {
		return err
	}

	err = client.Call("AdminHandler.Delete", admin, nil)
	if err != nil {
		return err
	}

	fmt.Printf("Admin \"%s\" deleted.\n", admin.Name)
	return nil
}

func reactivateAdmin(ctx *cli.Context) error {
	admin := admin.Admin{
		Name: ctx.String("name"),
	}

	client, err := getRPCClient()
	if err != nil {
		return err
	}

	err = client.Call("AdminHandler.Reactivate", admin, &admin)
	if err != nil {
		return err
	}

	fmt.Printf("Admin \"%s\" reactivated. All tokens before %s are now invalidated. New authorization token: %s\n", admin.Name, admin.ActivatedAt.String(), admin.AuthorizationToken)
	return nil
}

func setAdminExpiry(ctx *cli.Context) error {
	expiryString := ctx.String("expiry")
	var expiry time.Time
	if expiryString != "" {
		var err error
		expiry, err = time.Parse("2006-01-02", expiryString)
		if err != nil {
			return err
		}
	}

	admin := admin.Admin{
		Name:   ctx.String("name"),
		Expiry: expiry,
	}

	client, err := getRPCClient()
	if err != nil {
		return err
	}

	err = client.Call("AdminHandler.SetExpiry", admin, &admin)
	if err != nil {
		return err
	}

	fmt.Printf("Admin \"%s\" expiry set to %s.\n", admin.Name, admin.Expiry.String())
	return nil
}

func getRPCClient() (*rpc.Client, error) {
	return rpc.DialHTTP("tcp", "127.0.0.1:8081")
}
