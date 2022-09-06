package confirmer_interface

// moved in this package to avoid "import cycle not allowed"
import (
	"gitlab.altex.ro/ams/go_ddwcs/delivery/repository"
	logger "gitlab.altex.ro/plug/go_logger"
)

type ItemConfirm interface {
	IsConfirmed() bool
	ComputeConfirmation(stocks *map[int]int, isForcedConfirmation bool) bool
	GetWhRepo() *repository.WarehouseRepository
	GetConfirmedQty() int
	GetNeededQty() int // orderedQty
	GetRegionCode() string
	GetLogger() logger.Logger
	SaveState()
	RestoreState()
}
