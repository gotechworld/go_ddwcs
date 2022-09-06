package controller

import (
	"github.com/gorilla/mux"
	"gitlab.altex.ro/ams/go_ddwcs/awb/data_store"
	"gitlab.altex.ro/ams/go_ddwcs/awb/repository"
	"gitlab.altex.ro/ams/go_ddwcs/awb/service"
	"gitlab.altex.ro/ams/go_ddwcs/awb/service/adapter"
	"gitlab.altex.ro/ams/go_ddwcs/awb/util"
	lib "gitlab.altex.ro/plug/go_atx_lib"
	httpLocal "gitlab.altex.ro/plug/go_http"
	"net/http"
)

type StatusController struct{}

/**
 * Register the struct as a controller in the atx_lib
 * It will be used for the routing process
 */
func init() {
	lib.GetContainer().GetControllerRegistry().Register(NewStatusController())
}

func NewStatusController() *StatusController {
	return &StatusController{}
}

var statusActionWhitelistParams = []string{"increment_id", "customer_id||customer_email"}

/**
 * ReadStatusAction
 * @param r *http.Request
 * @return httpLocal.Response (atx_lib http.Response)
 * @desc - Controller action used to track the status of an awb
 */
func (ac *StatusController) ReadStatusAction(r *http.Request) httpLocal.Response {
	vars := mux.Vars(r)
	queryString := r.URL.Query()
	//user := r.Context().Value("user")

	logger := lib.GetContainer().GetLogger()
	controllerHelper := &util.StatusControllerHelper{}

	// validate required params - fail fast if something is missing
	if err := controllerHelper.ValidateParams(&queryString, statusActionWhitelistParams); err != nil {
		logger.Info("[StatusController]: " + err.Error())
		return httpLocal.NewErrorResponse("Valorile parametrior asteptati sunt gresite", http.StatusBadRequest)
	}

	// extract customer info, create cache key and get status from cache in case it exists
	customer := controllerHelper.ExtractCustomerFromQueryString(&queryString)
	cacheKey := queryString.Get("increment_id") + ":" + vars["awb"] + ":" + customer
	cacheAwbStatusStore := data_store.CreateCacheAwbStatusStoreInstance()
	if cachedStatus := cacheAwbStatusStore.Get(cacheKey); cachedStatus != nil {
		logger.Printf("[StatusController]: status loaded from cache for awb [%s]", vars["awb"])
		return httpLocal.NewResponse(cachedStatus, 0)
	}

	//check for awb existence
	orderAwbsProvider := repository.NewOrderAwbsRepo(queryString.Get("increment_id"), customer)
	validator := service.NewAwbValidator(orderAwbsProvider)
	if !validator.IsValid(vars["awb"]) {
		return httpLocal.NewErrorResponse("Awb-ul nu a fost gasit", http.StatusNotFound)
	}

	// get awb status
	statusTracker := &service.AwbStatusTracker{Logger: logger}
	status := statusTracker.GetAwbStatus(orderAwbsProvider.Get(vars["awb"]))

	// save details in cache
	go func(cacheKey string, status *adapter.AwbStatusDetails) {
		if status != nil {
			cacheAwbStatusStore.Save(cacheKey, status)
		}
	}(cacheKey, status)

	return httpLocal.NewResponse(status, 0)
}
