package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"rbac/cmd/internal"
	internaldomain "rbac/internal"
	"rbac/internal/elasticsearch"
	"rbac/internal/envvar"
	"rbac/internal/memcached"
	"rbac/internal/postgresql"
	"rbac/internal/service"

	"go.uber.org/zap"
)

func newProfile() internaldomain.Profile {
	return internaldomain.Profile{
		Profile_Picture:    "admin",
		Profile_Background: "admin",
		First_Name:         "admin",
		Last_Name:          "admin",
		Mobile:             "09123456789",
		Email:              "admin@admin.com",
	}
}
func newAccount() internaldomain.Account {
	return internaldomain.Account{
		UserName: "admin",
		Profile:  newProfile(),
	}
}
func newTask() []string {
	var tasks []string

	tasks = append(tasks, internaldomain.CREATE_ACCOUNT)
	tasks = append(tasks, internaldomain.GET_ACCOUNT)
	tasks = append(tasks, internaldomain.UPDATE_ACCOUNT)
	tasks = append(tasks, internaldomain.DELETE_ACCOUNT)
	tasks = append(tasks, internaldomain.LIST_ACCOUNT)

	tasks = append(tasks, internaldomain.CREATE_ROLE)
	tasks = append(tasks, internaldomain.GET_ROLE)
	tasks = append(tasks, internaldomain.UPDATE_ROLE)
	tasks = append(tasks, internaldomain.DELETE_ROLE)
	tasks = append(tasks, internaldomain.LIST_ROLE)

	tasks = append(tasks, internaldomain.CREATE_ACCOUNT_ROLE)
	tasks = append(tasks, internaldomain.GET_ACCOUNT_ROLE)
	tasks = append(tasks, internaldomain.UPDATE_ACCOUNT_ROLE)
	tasks = append(tasks, internaldomain.DELETE_ACCOUNT_ROLE)
	tasks = append(tasks, internaldomain.LIST_ACCOUNT_ROLE)

	tasks = append(tasks, internaldomain.CREATE_TASK)
	tasks = append(tasks, internaldomain.GET_TASK)
	tasks = append(tasks, internaldomain.UPDATE_TASK)
	tasks = append(tasks, internaldomain.DELETE_TASK)
	tasks = append(tasks, internaldomain.LIST_TASK)

	tasks = append(tasks, internaldomain.CREATE_ROLE_TASK)
	tasks = append(tasks, internaldomain.GET_ROLE_TASK)
	tasks = append(tasks, internaldomain.UPDATE_ROLE_TASK)
	tasks = append(tasks, internaldomain.DELETE_ROLE_TASK)
	tasks = append(tasks, internaldomain.LIST_ROLE_TASK)

	return tasks
}

func main() {
	var env string

	flag.StringVar(&env, "env", "", "Environment Variables filename")
	flag.Parse()

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(fmt.Errorf("new logger %w", err))
	}
	if err := envvar.Load(env); err != nil {
		log.Fatal(fmt.Errorf("envvar.Load %w", err))
	}
	vault, err := internal.NewVaultProvider()
	if err != nil {
		log.Fatal(fmt.Errorf("newVaultProvider %w", err))
	}
	conf := envvar.New(vault)
	db, err := internal.NewPostgreSQL(conf)
	if err != nil {
		log.Fatal(fmt.Errorf("newDB %w", err))
	}

	es, err := internal.NewElasticSearch(conf)
	if err != nil {
		log.Fatal(fmt.Errorf("new ElasticSearch %w", err))
	}
	m, err := internal.NewMemcached(conf)
	if err != nil {
		log.Fatal(fmt.Errorf("internal.NewMemcached %w", err))
	}

	token, err := internal.NewTokenMaker(conf)
	if err != nil {
		log.Fatal(fmt.Errorf("newTokenMaker %w", err))
	}
	repo := postgresql.NewRBAC(db)
	search := elasticsearch.NewRBAC(es)
	mclient := memcached.NewRBAC(m, search, logger)
	svc := service.NewRBAC(repo, mclient, token)

	//create new user
	ctx := context.Background()
	accId, err := svc.CreateAccount(ctx, newAccount(), "admin")
	if err != nil {
		log.Fatal(fmt.Errorf("new Account %w", err))
	}

	//create new role
	rid, err := svc.CreateRole(ctx, "ADMIN")
	if err != nil {
		log.Fatal(fmt.Errorf("new role %w", err))
	}

	//create new accountrole
	err = svc.CreateAccountRole(ctx, accId, rid)
	if err != nil {
		log.Fatal(fmt.Errorf("new accountrole %w", err))
	}

	for _, value := range newTask() {
		//create new task
		tid, err := svc.CreateTask(ctx, value)
		if err != nil {
			log.Fatal(fmt.Errorf("new task %w", err))
		}
		//create helptext for the task
		err = svc.CreateHelpText(ctx, internaldomain.HelpText{
			Task_id:  tid,
			HelpText: "Helptext " + value,
		})
		if err != nil {
			log.Fatal(fmt.Errorf("new helptext %w", err))
		}
		err = svc.CreateMenu(ctx, internaldomain.Menu{
			Task_id: tid,
			Name:    value,
		})
		if err != nil {
			log.Fatal(fmt.Errorf("new menu %w", err))
		}
		err = svc.CreateNavigation(ctx, internaldomain.Navigation{
			Task_id: tid,
			Name:    value,
		})
		if err != nil {
			log.Fatal(fmt.Errorf("new navigation %w", err))
		}
		//create new roletask
		err = svc.CreateRoleTask(ctx, tid, rid)
		if err != nil {
			log.Fatal(fmt.Errorf("new roletask %w", err))
		}
	}
}
