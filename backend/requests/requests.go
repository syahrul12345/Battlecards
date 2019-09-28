//Package requests package make all HTTP requests to a particular resource
package requests

import (
	"backend/utils"
	"io/ioutil"
	"net/http"
)

//Get returns the response when an API call is made
func Get(rpc string) map[string]interface{} {
	response, responseErr := http.Get(rpc)
	if responseErr != nil {
		return utils.Message(false, responseErr.Error())
	}
	//close the response to prevent overflows
	defer response.Body.Close()

	//read the response
	body, bodyErr := ioutil.ReadAll(response.Body)
	if bodyErr != nil {
		return utils.Message(false, bodyErr.Error())
	}
	return utils.Message(true, body)
}
