package tui
// ref: https://github.com/marcusolsson/tui-go/blob/master/example/chat/main.go

import (
	tui "github.com/marcusolsson/tui-go"
)

/**
version 1.0
*/

type LoginHandler func(string)

var loginChain *tui.SimpleFocusChain

type LoginView struct {
	tui.Box
	frame		 	*tui.Box
	loginHandler	LoginHandler

}

const logo = `     _____ __ ____  ___   ______________  
    / ___// //_/\ \/ / | / / ____/_  __/  
    \__ \/ ,<    \  /  |/ / __/   / /     
   ___/ / /| |   / / /|  / /___  / /      
  /____/_/ |_|  /_/_/ |_/_____/ /_/     `

func NewLoginView() *LoginView {
	view := &LoginView{}
	loginChain = &tui.SimpleFocusChain{}

	user := tui.NewEntry()
	user.SetFocused(true)
	user.SetSizePolicy(tui.Maximum, tui.Maximum)

	form := tui.NewGrid(0,0)
	form.AppendRow(tui.NewLabel("User"))
	form.AppendRow(user)

	login := tui.NewButton("[Login]")
	login.OnActivated(func(b *tui.Button) {
		if user.Text() != ""{
			if view.loginHandler != nil{
				view.loginHandler(user.Text())
			}
		}
	})
	buttons := tui.NewHBox(
		tui.NewSpacer(),
		tui.NewPadder(2,0,login),
	)

	window := tui.NewVBox(
		tui.NewPadder(10,1,tui.NewLabel(logo)),
		tui.NewPadder(12,0,tui.NewLabel("Welcome to Skynet 1.0!, Please login in first.")),
		tui.NewPadder(1,1,form),
		buttons,
	)
	window.SetBorder(true)

	wrapper := tui.NewVBox(
		tui.NewSpacer(),
		window,
		tui.NewSpacer(),
	)

	content := tui.NewHBox(tui.NewSpacer(), wrapper, tui.NewSpacer())

	root := tui.NewVBox(
		content,
	)
	loginChain.Set(user,login)

	view.frame = root
	view.Append(view.frame)

	return view
}

func (v *LoginView) OnLogin(handler LoginHandler) {
	v.loginHandler = handler
}