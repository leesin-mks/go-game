package com_game_component

type IComponent interface {
	/**
	 * 获取组件名称
	 *
	 * @return component simple name
	 */
	getName() string

	/**
	 * 初始化组件
	 *
	 * @return the init result
	 */
	initialize() bool

	/**
	 * 启动组件
	 *
	 * @return the start result
	 */
	// start() bool

	/**
	 * 停止组件
	 */
	stop()

	/**
	 * 重新加载组件
	 */
	reload() bool
}
