package presentation

import (
	"rogue/domain"

	"github.com/rthornton128/goncurses"
)

func Menu(r *Renderer) (int, string){
	r.window.Clear()
	r.window.MovePrint(10,10,"Enter 1 to Start New Game")
	r.window.MovePrint(11,10,"Enter 2 to Load Last Game")
	r.window.MovePrint(12,10,"Enter 3 to Watch Records")

	ch := r.window.GetChar()
	switch (ch){
	case '1':
		r.window.Clear()
		r.window.Refresh()
		r.window.MovePrint(10,10, "Write your name: ")
		goncurses.Echo(true)
		name, _ :=r.window.GetString(20);
		goncurses.Echo(false)
		return 1, name
	case '2':
		return 2, ""
	case '3':
		return 3, ""
	}
	

	return 0, ""
} 

func Records(r *Renderer, ScoreBoard *domain.LeaderBoards){
	r.window.Clear()
	r.window.Refresh()
	if len(ScoreBoard.Record)<1{
		r.window.MovePrint(1,10, "There is no records (Press q for menu)")
	} else {
		r.window.MovePrint(1,10, "Records (Press q to exit):")
	}
	for idx, elem := range ScoreBoard.Record{
		r.window.MovePrint(idx+2,10, elem.Name, ": ", elem.Record)
	}
	for {
		ch := r.window.GetChar()
		if ch == 'q'{
			break
		}
	}	
}