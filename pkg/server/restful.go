package server

type restfulRegisterer struct {
	s *server
}

func (rr *restfulRegisterer) Init() {
	rr.s.router.GET("kinds", rr.getKinds)
	rr.s.router.GET("resource/:kind", rr.getKeys)
	rr.s.router.GET("resource/:kind/*key", rr.getObject)

	rr.s.router.POST("eval", rr.eval)
	rr.s.router.POST("resource/:kind", rr.listObjects)
}
