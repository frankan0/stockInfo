package collectors

type CollectorManager struct {
	collectors []collector
}


func (m *CollectorManager) Start()  {
	m.collectors =  append(m.collectors,new(WeiboCollector))
	m.collectors =  append(m.collectors,new(WallCollector))
	//开始注册
	for _,c := range m.collectors{
		c.regTask()
	}
}






