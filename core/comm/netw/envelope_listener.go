package netw

type EnvelopeListener interface {
	OnNotify(notify *Notify)
	OnEvent(event *Event)
	OnStamp(stamp *Stamp)
	OnAddMoney(addMoney *AddMoney)
	OnDeal(deal *Deal)
	OnStand(stand *Stand)
	OnHit(hit *Hit)
	OnDouble(double *Double)
	OnPlayGame(playGame *PlayGame)
}
