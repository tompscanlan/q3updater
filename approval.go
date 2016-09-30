package q3updater

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	//	"time"
)

type Approval struct {
	Id          int    `json:"id,omitempty"`
	TeamId      int    `json:"teamID,omitempty"`
	Blob        int    `json:"blob,omitempty"`
	Description string `json:"description,omitempty"`
	Approved    bool   `json:"approved"`
}

var ApprovalId = 0

func NewApproval(blob []byte) *Approval {

	approval := new(Approval)
	approval.Blob = 0
	approval.TeamId = Team
	approval.Id = ApprovalId
	approval.Description = base64.StdEncoding.EncodeToString(blob)
	approval.Id += 1

	return approval
}

func GetApproved(server string, teamID int) ([]Approval, error) {
	var approvals []Approval

	// create approval request
	url := fmt.Sprintf("%s/%s?approved=%t&teamID=%d&limit=%d", server, "api/v1/approvables", true, Team, 10)
	log.Printf("GetApproved: %s", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("GetApproved", err)
		return approvals, err
	}
	req.Header.Set("Accept", "application/json")

	// make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("GetApproved", err)
		return approvals, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	if Verbose {
		log.Println("GetApproved: response Status:", resp.Status)
		log.Println("GetApproved: response Headers:", resp.Header)

		log.Println("GetApproved: response Body:", string(body))
	}
	err = json.Unmarshal(body, &approvals)
	if err != nil {
		log.Println("GetApproved", err)
		return approvals, err
	}
	return approvals, nil
}

func ParseApproved(jsonStr []byte) ([]Approval, error) {
	var approvals []Approval

	err := json.Unmarshal(jsonStr, approvals)
	if err != nil {
		log.Println(err)
		return approvals, err
	}

	return approvals, nil
}

func PostForApproval(server string, entry *JournalEntry) error {
	log.Println("PostForApproval: sending journal entry for approval")

	//	// add entry to the body
	//	jsonStr, err := json.Marshal(entry)
	//	if err != nil {
	//		log.Println(err)
	//		return err
	//	}

	decoded, err := base64.StdEncoding.DecodeString(entry.Message)
	if err != nil {
		log.Println(err)
	}

	approvable := NewApproval(decoded)
	approvableStr, err := json.Marshal(approvable)
	if err != nil {
		log.Println(err)
		return err
	}

	// create the request
	url := fmt.Sprintf("%s/%s", server, "api/v1/approvables")
	log.Printf("PostForApproval: %s, data: %s", url, approvableStr)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(approvableStr))
	if err != nil {
		log.Println(err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	// send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return err
	}

	defer resp.Body.Close()

	log.Println("PostForApproval: response Status:", resp.Status)
	log.Println("PostForApproval: response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println("PostForApproval: response Body:", string(body))
	return nil
}
