package nodes

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type NewNodeHandler struct {
	Path           string
	NodeRepository NodeRepository
}

func (h *NewNodeHandler) New(writer http.ResponseWriter, request *http.Request) {
	method := request.Method

	if method == "GET" {
		all, _ := h.NodeRepository.FindAll()
		writer.Header().Add("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)

		bytes, _ := json.Marshal(&all)
		_, _ = writer.Write(bytes)

		return
	}

	if method != "POST" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	requestBody, _ := ioutil.ReadAll(request.Body)
	nodeRequest := Node{}
	err := json.Unmarshal(requestBody, &nodeRequest)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.NodeRepository.Save(&nodeRequest)
	if err != nil {
		writer.WriteHeader(http.StatusConflict)
		return
	}

	writer.WriteHeader(201)
	writer.Header().Add("Content-Type", "application/json")
	bytes, _ := json.Marshal(&nodeRequest)
	_, _ = writer.Write(bytes)
}
