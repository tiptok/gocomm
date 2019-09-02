package config


//var (
//	once sync.Once
//)
//
//func Init(addr, name string) {
//	once.Do(func() {
//		source := grpc.NewSource(
//			grpc.WithAddress(addr),
//			grpc.WithPath(name),
//		)
//
//		if err := config.Load(source); err != nil {
//			log.Fatal(err)
//			return
//		}
//
//		go func() {
//			watcher, err := config.Watch()
//			if err != nil {
//				log.Fatal(err)
//			}
//
//			for {
//				v, err := watcher.Next()
//				if err != nil {
//					log.Log(err)
//					continue
//				}
//
//				log.Logf("[Init] file change: %v", string(v.Bytes()))
//			}
//		}()
//	})
//}
//
//func Get(conf interface{}, path ...string) (err error) {
//	if v := config.Get(path...); v != nil {
//		err = v.Scan(conf)
//	} else {
//		err = fmt.Errorf("[Get] 配置不存在, err: %v", path)
//	}
//	return
//}