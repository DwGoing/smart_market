package marketServcie

import "log"

// +ioc:autowire=true
// +ioc:autowire:type=singleton
// +ioc:autowire:constructFunc=NewMarketService

type MarketService struct {
	BinanceProvider *BinanceProvider `singleton:""`
}

func NewMarketService(service *MarketService) (*MarketService, error) {
	log.SetPrefix("[MarketService] ")
	return service, nil
}
