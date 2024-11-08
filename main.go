package main

import "api-gateway/bootstrap"

func main() {
	go func() {
		gatewayApp := bootstrap.NewGatewayApp()
		defer gatewayApp.Close()

		gatewayApp.Initialize()
		gatewayApp.SetupRoutes()
		gatewayApp.Run()
	}()

	go func() {
		managementApp := bootstrap.NewManagementApp()
		managementApp.Initialize()
		managementApp.SetupRoutes()
		managementApp.Run()
	}()

	select {}
}
