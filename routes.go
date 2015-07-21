package main

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
	Route{
		"EmulatorIndex",
		"GET",
		"/emulators",
		ShowEmulator,
	},
	Route{
		"EmulatorCreate",
		"POST",
		"/emulators",
		StartEmulator,
	},
	Route{
		"EmulatorDelete",
		"DELETE",
		"/emulators",
		StopEmulator,
	},
	Route{
		"Hubs",
		"GET",
		"/hubs",
		HubIndex,
	},
	Route{
		"HubCreate",
		"POST",
		"/hubs/{hubId}",
		HubCreate,
	},
	Route{
		"HubShow",
		"GET",
		"/hubs/{hubId}",
		HubShow,
	},
	Route{
		"HubDelete",
		"DELETE",
		"/hubs/{hubId}",
		HubDelete,
	},
	Route{
		"PortAttach",
		"POST",
		"/ports/{hubId}/{deviceId}",
		PortAttach,
	},
	Route{
		"PortDetach",
		"DELETE",
		"/ports/{hubId}/{deviceId}",
		PortDetach,
	},
}
