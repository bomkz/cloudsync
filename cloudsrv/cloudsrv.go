package cloudsrv

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/bomkz/cloudsync/cliauth"
	"github.com/bomkz/cloudsync/global"
)

func GetPilots() ([]string, error) {

	request := GetPilotsReqStruct{}

	request.Req = "getPilots"

	token, err := cliauth.GetAuthToken()
	if err != nil {
		return nil, err

	}
	request.Auth = token
	reqByte, err := json.Marshal(request)
	if err != nil {
		return nil, err

	}
	resp, err := http.Post(
		"http://localhost:9999/v1",
		"application/json",
		bytes.NewReader(reqByte),
	)
	if err != nil {
		return nil, err

	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	pilotsResp := global.GetPilotsResponseStruct{}

	if resp.StatusCode == 200 {
		json.Unmarshal(body, &pilotsResp)
		return pilotsResp.Body, nil
	}

	statusResponse := global.ResponseStruct{}
	err = json.Unmarshal(body, &statusResponse)
	if err != nil {
		return nil, errors.New(fmt.Sprint(resp.StatusCode) + ": " + err.Error())
	}
	return nil, errors.New(fmt.Sprint(resp.StatusCode) + ": " + statusResponse.Body)

}

func UpdatePilot(name string, pilot global.PilotsFile) error {

	request := global.UpdatePilotRequestStruct{}
	token, err := cliauth.GetAuthToken()
	if err != nil {
		return err
	}
	request.Auth = token
	request.Request = "updatePilot"
	request.Body.PilotData = pilot
	request.Body.Name = name

	reqBytes, err := json.Marshal(request)

	if err != nil {
		return err
	}
	resp, err := http.Post("http://localhost:9999/v1", "application/json", bytes.NewReader(reqBytes))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var pilotResp global.UpdatePilotRespStruct

	err = json.Unmarshal(body, &pilotResp)

	if err != nil {
		return err
	}

	switch pilotResp.Req {
	case "success":
		return nil
	case "error":
		return errors.New(pilotResp.Info)
	}
	return nil
}

func GetPilot(name, version string) (global.PilotsFile, error) {

	request := global.GetPilotReqStruct{}
	token, err := cliauth.GetAuthToken()
	if err != nil {
		return global.PilotsFile{}, err
	}
	request.Auth = token
	request.Pilot.Name = name
	parsedVersion, err := strconv.Atoi(version)
	if err != nil {
		return global.PilotsFile{}, fmt.Errorf("invalid pilot version %q: %w", version, err)
	}
	request.Pilot.Version = parsedVersion

	request.Req = "getPilot"

	reqByte, err := json.Marshal(request)

	if err != nil {
		return global.PilotsFile{}, err
	}

	resp, err := http.Post("http://localhost:9999/v1", "application/json", bytes.NewReader(reqByte))
	if err != nil {
		return global.PilotsFile{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return global.PilotsFile{}, err
	}

	pilotResp := global.GetPilotRespStruct{}

	err = json.Unmarshal(body, &pilotResp)
	if err != nil {
		return global.PilotsFile{}, err
	}

	return pilotResp.Pilot, nil
}

func GetPilotInfo(name string) (int, error) {
	request := global.GetPilotInfoReqStruct{}

	request.Req = "getPilotInfo"
	request.Name = name

	token, err := cliauth.GetAuthToken()
	if err != nil {
		return 0, err

	}
	request.Auth = token
	reqByte, err := json.Marshal(request)
	if err != nil {
		return 0, err

	}
	resp, err := http.Post(
		"http://localhost:9999/v1",
		"application/json",
		bytes.NewReader(reqByte),
	)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	pilotsResp := global.GetPilotInfoRespStruct{}

	if resp.StatusCode == 200 {
		json.Unmarshal(body, &pilotsResp)
		return pilotsResp.Content, nil
	}

	statusResponse := global.ResponseStruct{}
	err = json.Unmarshal(body, &statusResponse)
	if err != nil {
		return 0, errors.New(fmt.Sprint(resp.StatusCode) + ": " + err.Error())
	}
	return 0, errors.New(fmt.Sprint(resp.StatusCode) + ": " + statusResponse.Body)

}
