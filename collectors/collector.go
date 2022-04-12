package collectors

type collector interface {
	regTask()
	start()
}
