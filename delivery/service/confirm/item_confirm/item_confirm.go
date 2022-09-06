package item_confirm

// moved in this package to avoid "import cycle not allowed"
import (
	"gitlab.altex.ro/ams/go_ddwcs/delivery/repository"
	logger "gitlab.altex.ro/plug/go_logger"
)

type ItemConfirm struct {
	WhRepo                 *repository.WarehouseRepository
	Logger                 logger.Logger
	Confirmations          map[int]int //map[warehouseId: confirmedQty, ..]
	ForcedConfirmations    map[int]int //map[warehouseId: confirmedQty, ...}
	ConfirmedQty           int
	NeededQty              int // orderedQty
	HasForcedConfirmations bool
	RegionCode             string
	itemConfirmSavePoints  []ItemConfirm
}

// ItemConfirm - constructor
func NewItemConfirm(whRepo *repository.WarehouseRepository, logger logger.Logger, neededQty int, deliveryRegion string) *ItemConfirm {
	return &ItemConfirm{
		WhRepo:                whRepo,
		Logger:                logger,
		Confirmations:         map[int]int{},
		ForcedConfirmations:   map[int]int{},
		NeededQty:             neededQty,
		itemConfirmSavePoints: []ItemConfirm{},
		RegionCode:            deliveryRegion,
	}
}

// Check if current item is confirmed.
//
// @return bool
//
func (c ItemConfirm) IsConfirmed() bool {
	if c.ConfirmedQty == c.NeededQty {
		return true
	}

	return false
}

// Compute confirmation for current item.
//
// @param stocks map[int]int - map[warehouseId] = warehouseQty
// @param isForcedConfirmation bool - Flag for warehouses pool - if we should enforce qty
// @return bool
//
func (c *ItemConfirm) ComputeConfirmation(stocks *map[int]int, isForcedConfirmation bool) bool {
	for whId, qty := range *stocks {
		if c.IsConfirmed() {
			return true
		}
		// we already confirmed from this warehouse
		if 0 != c.Confirmations[whId] {
			return false
		}

		// we assure that $qty <= ($this->neededQty - $this->confirmedQty)
		if c.ConfirmedQty+qty > c.NeededQty {
			qty = c.NeededQty - c.ConfirmedQty
		}

		// for ForcedConfirmation we can confirm smaller qty's
		if isForcedConfirmation {
			c.HasForcedConfirmations = true
			c.ForcedConfirmations[whId] = qty
			c.Confirmations[whId] = qty
			c.ConfirmedQty += qty

			if c.IsConfirmed() {
				break
			}

			continue
		}

		// for not ForcedConfirmation we should confirm only when qty's match
		if c.ConfirmedQty+qty >= c.NeededQty {
			c.Confirmations[whId] = qty
			c.ConfirmedQty += qty
		}

	} // for stocks

	return c.IsConfirmed()
}

// Save current state for later restoreState() in a stack manner.
func (c *ItemConfirm) SaveState() {
	c.itemConfirmSavePoints = append(c.itemConfirmSavePoints, *c)
}

// Restore state to the last saveState() in a stack manner.
func (c *ItemConfirm) RestoreState() {
	if 0 > len(c.itemConfirmSavePoints)-1 {
		return
	}

	itemConfirmSavePoints := c.itemConfirmSavePoints[len(c.itemConfirmSavePoints)-1]
	c.itemConfirmSavePoints = c.itemConfirmSavePoints[:len(c.itemConfirmSavePoints)-1]
	c.Confirmations = itemConfirmSavePoints.Confirmations
	c.ForcedConfirmations = itemConfirmSavePoints.ForcedConfirmations
	c.ConfirmedQty = itemConfirmSavePoints.ConfirmedQty
	c.NeededQty = itemConfirmSavePoints.NeededQty
	c.HasForcedConfirmations = itemConfirmSavePoints.HasForcedConfirmations
	c.RegionCode = itemConfirmSavePoints.RegionCode
}

func (c ItemConfirm) GetWhRepo() *repository.WarehouseRepository {
	return c.WhRepo
}

func (c ItemConfirm) GetConfirmedQty() int {
	return c.ConfirmedQty
}

func (c ItemConfirm) GetNeededQty() int {
	return c.NeededQty
}

func (c ItemConfirm) GetRegionCode() string {
	return c.RegionCode
}

func (c ItemConfirm) GetLogger() logger.Logger {
	return c.Logger
}

func (c ItemConfirm) GetConfirmation() map[int]int {
	return c.Confirmations
}
