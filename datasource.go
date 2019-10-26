package echor

import (
	"errors"
	"sync"
)

// DbDriverInterface class
type DbDriverInterface interface {
	GetNativeConnector() interface{}
}

// Source struct
type Source struct {
	Driver   string
	Address  string
	Port     int
	Username string
	Password string
	DbIndex  int
	DbName   string
}

func createDataSourceDriver(source Source) (DbDriverInterface, error) {
	if source.Driver == "mysql" {
		driver := NewTsqlDriver(source)
		return driver, nil
	} else if source.Driver == "redis" {
		driver := NewRedisDriver(source)
		return driver, nil
	}
	return nil, errors.New("DATASOURCE_DRIVER_NOT_SUPPORTED")
}

// Connector struct
type Connector struct {
	source Source
	driver DbDriverInterface
}

// Datasources struct
type Datasources struct {
	connectors map[string]Connector
}

var datasourcesInstance *Datasources

var once sync.Once

func instance() *Datasources {
	once.Do(func() {
		datasourcesInstance = &Datasources{
			connectors: make(map[string]Connector),
		}
	})
	return datasourcesInstance
}

// Register func
func Register(datasourceID string, source Source) (err error) {
	datasources := instance()
	if _, ok := datasources.connectors[datasourceID]; ok {
		err = errors.New("DUP_DATASOURCE_ID")
		return
	}
	driver, err := createDataSourceDriver(source)
	datasources.connectors[datasourceID] = Connector{
		source: source,
		driver: driver,
	}
	return err
}

func getDatasourceDriver(datasourceID string, driver string) interface{} {
	datasources := instance()
	// fmt.Println(datasourceID, datasources.connectors)
	if connector, ok := datasources.connectors[datasourceID]; ok {
		if driver == "" || driver == connector.source.Driver {
			return connector.driver.GetNativeConnector()
		}
	}
	return nil
}
