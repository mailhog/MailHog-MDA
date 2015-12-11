package main

import (
	"flag"
	"log"
	"sync"

	"github.com/mailhog/MailHog-MDA/config"
	sconfig "github.com/mailhog/backends/config"
	"github.com/mailhog/backends/delivery"
	"github.com/mailhog/backends/mailbox"
	"github.com/mailhog/backends/resolver"
)

var deliveryService delivery.Service
var mailboxService mailbox.Service
var resolverService resolver.Service

var conf *config.Config
var wg sync.WaitGroup

func configure() {
	config.RegisterFlags()
	flag.Parse()
	conf = config.Configure()

	configureBackends()
}

func main() {
	configure()

	deliveries := make(chan *delivery.Message, 100)
	deliveryService.Deliveries(deliveries)

	for m := range deliveries {
		log.Printf("Received: %+v", m)
		/*
		   TODO
		   - split message (group by server etc)
		   - deliver local messages
		   - deliver remote messages (or queue for external delivery?)
		*/
		for _, t := range m.To {
			mbox, err := mailboxService.Open(t)
			if err != nil {
				log.Printf("Error opening mailbox %s: %s", t, err)
				continue
			}
			err = mbox.Store(*m)
			if err != nil {
				log.Printf("Error storing message for %s: %s", t, err)
			}
		}
		err := deliveryService.Delivered(*m, true)
		if err != nil {
			log.Printf("Error: %+v", m)
		}
	}
}

func configureBackends() error {
	var d, m, r sconfig.BackendConfig
	var err error

	if conf.Delivery != nil {
		log.Printf("got delivery config")
		d, err = conf.Delivery.Resolve(conf.Backends)
		if err != nil {
			return err
		}
	}

	if conf.Mailbox != nil {
		log.Printf("got mailbox config")
		m, err = conf.Mailbox.Resolve(conf.Backends)
		if err != nil {
			return err
		}
	}

	if conf.Resolver != nil {
		log.Printf("got resolver config")
		r, err = conf.Resolver.Resolve(conf.Backends)
		if err != nil {
			return err
		}
	}

	deliveryService = delivery.Load(d, *conf)
	resolverService = resolver.Load(r, *conf)
	mailboxService = mailbox.Load(m, *conf, resolverService)

	return nil
}
