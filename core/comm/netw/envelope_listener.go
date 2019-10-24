package netw

type EnvelopeListener interface {
	OnNotify(notify *Notify)
	OnEvent(c interface{}, event *Event)
	OnStamp(c interface{}, stamp *Stamp)
	OnAddMoney(c interface{}, addMoney *AddMoney)
	OnDeal(c interface{}, deal *Deal)
	OnStand(c interface{}, stand *Stand)
	OnHit(c interface{}, hit *Hit)
	OnDouble(c interface{}, double *Double)
	OnPlayGame(c interface{}, playGame *PlayGame)
}
