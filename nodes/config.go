package nodes

import (
	"fmt"

	"github.com/denkhaus/yamlconfig"
)

type NodesConfig struct {
	*yamlconfig.Config
}

///////////////////////////////////////////////////////////////////////////////////////////////
func (c *NodesConfig) GetLogglyTokenById(tokenId string) (string, error) {
	conf := c.GetObject("logging:loggly:tokens")
	for _, entry := range conf.([]interface{}) {
		inf := entry.(map[interface{}]interface{})
		if inf["name"] == tokenId {
			return inf["token"].(string), nil
		}
	}

	return EMPTY_STRING, fmt.Errorf("Configuration Error:: Loggly token for id %s not found.", tokenId)
}

///////////////////////////////////////////////////////////////////////////////////////////////
// GetCassandraConfig
///////////////////////////////////////////////////////////////////////////////////////////////
func (c *NodesConfig) GetMongoConfig() (map[string]interface{}, error) {
	conf := make(map[string]interface{})
	conf["host"] = c.GetString("mongo:host")
	return conf, nil
}

/////////////////////////////////////////////////////////////////////////////////////////////////
//
/////////////////////////////////////////////////////////////////////////////////////////////////
func (c *NodesConfig) GetStathatConfig() (map[string]interface{}, error) {
	conf := make(map[string]interface{})
	if userkey := c.GetString("metrics:stathat:userkey"); userkey == EMPTY_STRING {
		return nil, fmt.Errorf("Configuration Error:: No Stathat userkey found.")
	} else {
		conf["userkey"] = userkey
	}

	if duration := c.GetDuration("metrics:stathat:duration"); duration == DURATION_DAY {
		return nil, fmt.Errorf("Configuration Error:: No Stathat duration found.")
	} else {
		conf["duration"] = duration
	}

	return conf, nil
}

/////////////////////////////////////////////////////////////////////////////////////////////////
func (c *NodesConfig) GetLibratoConfig() (map[string]interface{}, error) {
	conf := make(map[string]interface{})
	if apitoken := c.GetString("metrics:librato:apitoken"); apitoken == EMPTY_STRING {
		return nil, fmt.Errorf("Configuration Error:: No Librato apitoken found.")
	} else {
		conf["apitoken"] = apitoken
	}

	if hostname := c.GetString("metrics:librato:hostname"); hostname == EMPTY_STRING {
		return nil, fmt.Errorf("Configuration Error:: No Librato hostname found.")
	} else {
		conf["hostname"] = hostname
	}

	if email := c.GetString("metrics:librato:email"); email == EMPTY_STRING {
		return nil, fmt.Errorf("Configuration Error:: No Librato email found.")
	} else {
		conf["email"] = email
	}

	if duration := c.GetDuration("metrics:librato:duration"); duration == DURATION_DAY {
		return nil, fmt.Errorf("Configuration Error:: No Librato duration found.")
	} else {
		conf["duration"] = duration
	}

	return conf, nil
}

/////////////////////////////////////////////////////////////////////////////////////////////////
//
/////////////////////////////////////////////////////////////////////////////////////////////////
func (c *NodesConfig) Init(fileName string) error {
	c.Config = yamlconfig.NewConfig(fileName)
	if err := c.Load(func(config *yamlconfig.Config) {
		/*
			   config.SetDefault("bsclient:hosts:host", "127.0.0.1")
			   config.SetDefault("bsclient:rpc:port", 5680)
			   config.SetDefault("bsclient:rpc:path", "/rpc")
			   config.SetDefault("bsclient:rpc:username", "your-username")
			   config.SetDefault("bsclient:rpc:password", "your-password")
			   config.SetDefault("bsclient:rpc:use-ssl", false)


			config.SetDefault("bsclient:default", "please specify defaultacct")
			config.SetDefault("bsclient:tip:account", "happytip")
			config.SetDefault("bsclient:tip:message", "From happyshares with love")

			config.SetDefault("storage:poolsize", 10)
			config.SetDefault("storage:network", "tcp")
			config.SetDefault("storage:address", ":6379")
			config.SetDefault("storage:password", "")
			config.SetDefault("storage:queuedb", 9)

			config.SetDefault("system:payout:minpayoutamount", 0.5)
			config.SetDefault("system:payout:multiplier", 0.5)
			config.SetDefault("system:payout:payoutaccount", "!Specify PayoutAccount")
			config.SetDefault("system:payout:trxcomment", "From happyshares with love")
		*/
	}, "", false); err != nil {
		return err
	}

	return nil
}

/////////////////////////////////////////////////////////////////////////////////////////////////
//
/////////////////////////////////////////////////////////////////////////////////////////////////
func CreateNodesConfig(fileName string) (*NodesConfig, error) {
	cfig := &NodesConfig{}
	if err := cfig.Init(fileName); err != nil {
		return nil, err
	}
	return cfig, nil
}
