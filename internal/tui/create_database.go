package tui

import (
	"citadel/internal/api"
	"citadel/internal/cli"
	"citadel/internal/util"
	"fmt"
	"os"
	"strconv"

	"github.com/aquasecurity/table"
	"github.com/sveltinio/prompti/input"
)

func CreateDatabase(organizationSlug string, projectSlug string) {
	dbmsModel := newChooseDBMS()

	dbms, err := dbmsModel.Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	databaseName, err := input.Run(
		&input.Config{
			Message:      "What's the name of your database?",
			Placeholder:  "webapp-db",
			ValidateFunc: util.SlugValidateFunc,
		},
	)
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	var databaseUsername string
	if dbms.ID == "postgres" || dbms.ID == "mysql" {
		databaseUsername, err = input.Run(
			&input.Config{
				Message:     "What's the username of your database?",
				Placeholder: "steve",
			},
		)
		if err != nil {
			fmt.Println("Error running program:", err)
			os.Exit(1)
		}
	}

	databasePassword, err := input.Run(
		&input.Config{
			Message: `What's the password of your database?`,
		},
	)
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	diskSize, err := input.Run(
		&input.Config{
			Message:     "How much disk space do you need for your database in GB, between 1 and 100 (volumes are billed at $0.50/GB/month)?",
			Placeholder: "10",
			ValidateFunc: func(s string) error {
				if s == "" {
					return nil
				}
				if _, err := strconv.Atoi(s); err != nil {
					return err
				}
				return nil
			},
		},
	)
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	diskSizeInt, err := strconv.Atoi(diskSize)
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	pricePerMonth := float64(diskSizeInt)*0.5 + 6.38

	fmt.Println("")

	shouldCreate := cli.AskYesOrNo("This will create a " + diskSize + "GB " + dbms.Name + " database for $" + fmt.Sprintf("%.2f", pricePerMonth) + "/month. Do you want to continue?")
	if !shouldCreate {
		return
	}

	connectionString, databaseSlug, err := api.CreateDatabase(
		organizationSlug,
		projectSlug,
		dbms.ID,
		databaseName,
		databaseUsername,
		databasePassword,
		diskSize,
	)
	if err != nil {
		fmt.Println("Error creating database:", err)
		os.Exit(1)
	}

	fmt.Println("✅\nDatabase created successfully!")

	t := table.New(os.Stdout)

	t.SetHeaders("", "Database Credentials")
	t.SetHeaderStyle(table.StyleBold)
	t.SetAutoMergeHeaders(true)

	t.AddRow("Host", "citadel-database-"+databaseSlug+".internal")
	t.AddRow("Name", databaseSlug)
	if dbms.ID == "postgres" || dbms.ID == "mysql" {
		t.AddRow("Username", databaseUsername)
	}
	t.AddRow("Password", databasePassword)
	t.AddRow("Disk size", diskSize+" GB")

	t.Render()

	fmt.Println("Connection string:", connectionString)
}

func newChooseDBMS() SelectModel {
	choices := []SelectChoice{
		{
			Name: "PostgreSQL",
			ID:   "postgres",
		},
		{
			Name: "MySQL",
			ID:   "mysql",
		},
		{
			Name: "Redis",
			ID:   "redis",
		},
	}

	return NewSelectModel("Choose a database", choices)
}
