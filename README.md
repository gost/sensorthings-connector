# sensorthings-connector
Service for fetching sensor data from different services and sending it to a SensorThings server. The connector is plugin based which makes it easy to create your own modules, another usefull case is when you have multiple different API accounts for fetching certain sensor data, just copy the module, rename it and set it up using a config file.

## Configuration
On startup the connector looks for config.json if the -config flag is not supplied, the file contains the following  

```
{
    // connector config
    "connector":{
      "host":"0.0.0.0", // string (ip to run the HTTP server on)
      "port":5000, // int (port to run HTTP server on)
      "modulePath": "", // path to module folder leave empty to use program location (os.Args[0])
      "startModulesOnStartup": true // bool (start the modules on startup, if set to false modules must be started using the REST service)
    },
    // logging config
    "logging": {
      "status": {
        "enabled": false, // bool (set to true to log a status report every xx seconds)
        "intervalSeconds": 3600 // int (how much seconds between status reports)
      },
      "discord": {
        "enabled": false, // bool (set to true to send logs from level info and up to a Discord server)
        "name": "Connector", // string (name to report the status in Discord)
        "hook": "" // string (discord hook url, this url can be generated as admin of the Discord server)
      }
    }
}
```

## Logging
The connector logs to Stderr and can also be setup to log to Discord, just set it up using config.json. It is also possible to create a status report for a time interval, this can also be enabled using config.json 

## REST service
The connector contains a HTTP server for various purposes

### GET /Modules
To see the current loaded modules and their status browse to host:port/Modules

### POST /Modules/State
A module can be started/stopped from the /Modules/State endpoint  


POST body
```
{
    "on": true, // bool (true = start module, false = stop module)
    "moduleId": "" // string (id of the module, set by the module config or auto generated, check /Modules endpoint for id)
}
```

Status 400 when sending incorrect body or module not found  
Status 500 when the module is crashed and cannot be started  
Status 200 if no problem occured and the module started/stopped  

Response body contains errors explaining the error when status 400 or 500 was send back 

### /moduleid/xxx
Every module can expose their own endpoints to see which endpoints are available for a module check out /Modules

## Modules (Plugins)
You can write your own modules by using ConnectorModuleBase for examples check modules/netatmo or modules/foobot  

At startup the connector will search for plugin files which end with .so and tries to load them as a connecor module. Currently building and running plugins is only supported on Linux, to build a plugin run the following

```
$ cd modules/netatmo
$ go build -buildmode=plugin -o netatmo.so main.go
```

The current modules Netatmo and Foobot expect a .json config file in the same directory with the same name as the module. For example when netatmo1.so is loaded it tries to load netatmo1.json from the same directory.  

## ConnectorModuleBase
ToDo
