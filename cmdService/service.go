package cmdService

import (
	"fmt"
	"time"

	"github.com/MickMake/GoUnify/Only"
	"github.com/MickMake/GoUnify/cmdLog"
	"github.com/kardianos/service"
)


func (c *Service) Control(action string) error {
	for range Only.Once {

		// errs := make(chan error, 5)
		// logger, err = s.Logger(errs)
		// if err != nil {
		// 	break
		// }
		//
		// go func() {
		// 	for {
		// 		err := <-errs
		// 		if err != nil {
		// 			log.Print(err)
		// 		}
		// 	}
		// }()

		if action == "" {
			fmt.Printf("Valid actions: %q\n", service.ControlAction)
			break
		}

		c.Error = service.Control(c.service, action)
		if c.Error != nil {
			fmt.Printf("Valid actions: %q\n", service.ControlAction)
			break
		}
		fmt.Printf("Service action '%s' OK.\n", action)

		// err = s.Run()
		// if err != nil {
		// 	break
		// }
	}

	return c.Error
}

// func ServiceStart() error {
// 	var err error
//
// 	for range Only.Once {
//
// 	}
//
// 	return err
// }
//
// func ServiceStop() error {
// 	var err error
//
// 	for range Only.Once {
//
// 	}
//
// 	return err
// }
//
// func ServiceState() error {
// 	var err error
//
// 	for range Only.Once {
//
// 	}
//
// 	return err
// }


func (p *program) Start(s service.Service) error {
	if service.Interactive() {
		cmdLog.Printf("Running in terminal.")
	} else {
		cmdLog.Printf("Running under service manager.")
	}

	p.exit = make(chan struct{})

	// Start should not block. Do the actual work async.
	go p.exec()
	return nil
}

func (p *program) exec() {
	cmdLog.Printf("I'm running %v.", service.Platform())

	cmdLog.Printf("app.StartApp()")
	time.Sleep(time.Second * 120)

	ticker := time.NewTicker(2 * time.Second)
	for {
		select {
			case tm := <-ticker.C:
				cmdLog.Printf("Still running at %v...", tm)
			case <-p.exit:
				ticker.Stop()
				return
		}
	}
}

func (p *program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	cmdLog.Printf("I'm Stopping!")
	close(p.exit)
	return nil
}

