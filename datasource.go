package echor

import (
	"errors"
	"sync"
)

// DbDriverInterface class
type DbDriverInterface interface {
	GetNativeConnector() interface{}
}

// DatasourceSpecs struct
type DatasourceSpecs struct {
	Driver   string
	Address  string
	Port     int
	Username string
	Password string
	DbIndex  int
	DbName   string
}

func createDataSourceDriver(source DatasourceSpecs) (DbDriverInterface, error) {
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
type datasourceConnector struct {
	source DatasourceSpecs
	driver DbDriverInterface
}

// Datasources struct
type Datasources struct {
	connectors map[string]datasourceConnector
}

var datasourcesInstance *Datasources

var once sync.Once

func instance() *Datasources {
	once.Do(func() {
		datasourcesInstance = &Datasources{
			connectors: make(map[string]datasourceConnector),
		}
	})
	return datasourcesInstance
}

// RegisterDatasource func
func RegisterDatasource(datasourceID string, source DatasourceSpecs) (err error) {
	datasources := instance()
	if _, ok := datasources.connectors[datasourceID]; ok {
		err = errors.New("DUP_DATASOURCE_ID")
		return
	}
	driver, err := createDataSourceDriver(source)
	datasources.connectors[datasourceID] = datasourceConnector{
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
