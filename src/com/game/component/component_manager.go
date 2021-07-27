package com_game_component

var modules = make(map[string]IComponent)

func addComponent(name string) bool {

	return true
}

func start() bool {
	if initModules() {
		return false
	}
	return startModules()
}

func initModules() bool {

	return true
}

func startModules() bool {

	return true
}
