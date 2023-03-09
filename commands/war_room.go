package commands

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
)

// default google meet links
var serviceData = map[string]string{
	"vrf":     "https://meet.google.com/bcc-cqid-okb",
	"feeds":   "https://meet.google.com/jsg-cvtz-sxe",
	"keepers": "https://meet.google.com/nqa-mavr-qsr",
	"other":   "https://meet.google.com/hqh-vcbr-mtx",
}

type PdDetails struct {
	Summary           string
	AlertName         string
	NetworkName       string
	Job               string
	RelevantGMeetLink string
	PdLink            string
	CllService        string

	ExplorerDetailKey string
	Contract          string
	ContractAddress   string
	RegistryAddress   string
	KeyHash           string
	KeyUnderscoreHash string
}

func PostSlackWarRoomMessage(detailsByLines []string, pdLink string, slackEndpoint string) (int, error) {
	// try best to parse alert details
	var details PdDetails
	details.Init(detailsByLines, pdLink)

	// build message.
	msg := MakeSlackMsgDetails(details)

	// post message to Slack
	m := bytes.NewReader([]byte(msg))
	resp, err := http.Post(slackEndpoint, "application/json", m)
	if err != nil {
		fmt.Print("Failed post request to Slack\n")
		return resp.StatusCode, err
	}

	return resp.StatusCode, nil
}

func (details *PdDetails) Init(detailsByLines []string, pdLink string) error {
	// set PagerDuty link
	details.PdLink = pdLink

	// get alert summary; most important.
	for _, line := range detailsByLines {
		if strings.Contains(line, "- summary") {
			details.Summary = strings.Replace(strings.Split(line, "- summary = ")[1], `"`, `'`, -1)
			break
		}
	}

	if details.Summary == "" {
		fmt.Print("Missing alert summary when searching for label '- summary'. Using first line in file!")

		strangeAlert := ""
		for _, line := range detailsByLines {
			strangeAlert += line + " "
		}

		details.Summary = strangeAlert
	}

	details.Summary = strings.Replace(details.Summary, `<`, `less than`, -1)
	details.Summary = strings.Replace(details.Summary, `>`, `greater than`, -1)
	details.Summary = strings.Replace(details.Summary, `&`, `&amp;`, -1)
	details.Summary = strings.Replace(details.Summary, `"`, `\"`, -1)
	details.Summary = strings.Replace(details.Summary, `*`, `\*`, -1)

	// parse rest of details from file input (alertname, relevant explorer information, network name)
	for _, line := range detailsByLines {
		if strings.Contains(line, "- alertname") {
			details.AlertName = strings.Split(line, "- alertname = ")[1]

		} else if strings.Contains(line, "- network_name = ") {
			details.NetworkName = strings.Split(line, "- network_name = ")[1]

		} else if strings.Contains(line, "- job =") {
			if strings.Contains(line, "ocr_telemetry_prometheus_exporter_") {
				details.NetworkName = strings.Split(line, "ocr_telemetry_prometheus_exporter_")[1]
			}
			if strings.Contains(line, "atlas_prod_otpe2_") {
				details.NetworkName = strings.Split(line, "atlas_prod_otpe2_")[1]
			}

		} else if strings.Contains(line, "- contract =") {
			details.Contract = strings.Split(line, "- contract = ")[1]

		} else if strings.Contains(line, "- contract_address = ") {
			details.ContractAddress = strings.Split(line, "- contract_address = ")[1]

		} else if strings.Contains(line, "- key_hash = ") {
			details.KeyHash = strings.Split(line, "- key_hash = ")[1]

		} else if strings.Contains(line, "- registry_address = ") {
			details.RegistryAddress = strings.Split(line, "- registry_address = ")[1]

		}
	}

	// get details based on AlertName
	if details.AlertName != "" {
		// get correct GoogleMeet link based off service name
		details.RelevantGMeetLink = serviceData[strings.ToLower(details.CllService)]

		// do best to gather more information using input:
		//
		//	guess service name based on alert name
		if strings.Contains(strings.ToLower(details.AlertName), "feed") ||
			strings.Contains(strings.ToLower(details.AlertName), "offchainaggregatoranswerstalled") ||
			strings.Contains(strings.ToLower(details.AlertName), "highoffchainaggregatorexpectedanswervsonchainanswerdeviation") ||
			strings.Contains(strings.ToLower(details.AlertName), "consensusfailurewarning") {
			details.CllService = "feeds"
		} else if strings.Contains(strings.ToLower(details.AlertName), "vrf") {
			details.CllService = "vrf"
		} else if strings.Contains(strings.ToLower(details.AlertName), "upkeep") {
			details.CllService = "keepers"
		}
	}

	// all caps the name "VRF" if is vrf
	details.CllService = strings.Title(strings.ToLower(details.CllService))
	if details.CllService == "Vrf" {
		details.CllService = "VRF"
	}

	// fill if unfound
	if details.RelevantGMeetLink == "" {
		details.RelevantGMeetLink = serviceData["other"]
	}

	return nil
}

func GetPdLinkFromUser() string {
	// get PD link:
	goodPdLink := false
	var userInput string

	for !goodPdLink {
		fmt.Println("Paste in PagerDuty alert link: ")
		fmt.Scanln(&userInput)

		if IsUrl(userInput) {
			goodPdLink = true
			return userInput
		} else {
			fmt.Print("Received a bad link. Try again. Input: " + userInput)
		}
	}

	return "https://chainlink.pagerduty.com/service-directory/PFO8X1C"
}

func MakeSlackMsgDetails(details PdDetails) string {
	explorerDetailKey, explorerDetail := getExplorerDetailsDependingOnProduct(details)

	// TO DO: make struct to marshal into JSON
	msg := `{
		"text": "New war room: ` + details.Summary + `",
		"blocks": [
			{
				"type": "header",
				"text": {
					"type": "plain_text",
					"text": "A new war-room has been created!",
					"emoji": true
				}
			},
			{
				"type": "section",
				"text": {
					"type": "mrkdwn",
					"text": "<` + details.PdLink + `|*` + details.Summary + `*>`

	if details.CllService != "" {
		msg += `\n\n*Service* \n` + details.CllService
	}

	if details.AlertName != "" || details.NetworkName != "" || explorerDetailKey != "" {
		msg += `\n\n *Details*\n`
	}

	if details.AlertName != "" {
		msg += "\\n`alertname` = " + details.AlertName
	}

	if details.NetworkName != "" {
		msg += "\\n`network` = " + details.NetworkName
	}

	if explorerDetailKey != "" {
		msg += "\\n`" + explorerDetailKey + "` = " + explorerDetail
	}

	msg += `"

				}
			},
			{
				"type": "divider"
			},
			{
				"type": "section",
				"text": {
					"type": "mrkdwn",
					"text": "*Communications*\n:google-meet-intensifies:<` + details.RelevantGMeetLink + `|Google Meets>"
				}
			}
		]
	}
	`

	return msg
}

func MakeGenMsgDetails(details PdDetails) string {
	explorerDetailKey, explorerDetail := getExplorerDetailsDependingOnProduct(details)

	// build message:
	msg := "\nThere is a new incident!"
	msg += "\n\n" + details.Summary + "\n" + details.PdLink

	if details.CllService != "" {
		msg += "\n\nService: " + details.CllService
	}

	if details.AlertName != "" || details.NetworkName != "" || explorerDetailKey != "" {
		msg += "\n\n" + "Details:"
	}

	if details.AlertName != "" {
		msg += "\n- alertName = " + details.AlertName
	}

	if details.NetworkName != "" {
		msg += "\n- networkName = " + details.NetworkName
	}

	if explorerDetailKey != "" {
		msg += "\n- " + explorerDetailKey + " = " + explorerDetail
	}

	msg += "\n\nCommunications:"
	msg += "\nGoogle Meets: " + details.RelevantGMeetLink

	return msg
}

func MakeFileForUserForInput(filepath string) (*os.File, error) {
	// creates a file, opens it for the user to input PagerDuty details, closes, and returns the file.

	fmt.Println("Creating file; paste in PagerDuty details and then save. Press enter to continue.")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()

	// create file
	filePdDetails, err := os.Create(filepath)
	if err != nil {
		fmt.Print(err)
		return filePdDetails, err
	}
	filePdDetails.Close()

	// default open file for mac os
	err = exec.Command("open", filePdDetails.Name()).Run()
	if err != nil {
		// if err, try for linux directly with gvim, required.
		err = exec.Command("gvim", filePdDetails.Name()).Run()
		if err != nil {
			fmt.Printf(err.Error(), "Unable to open new file ")
			return filePdDetails, err
		}
	}

	// wait until saved
	fmt.Println("Press Enter when done saving.")
	input.Scan()

	return filePdDetails, nil
}

func GetListFromFileLines(file *os.File) ([]string, error) {
	// build string list of file lines
	var lines []string

	file, err := os.Open(file.Name())
	if err != nil {
		return lines, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return lines, err
	}
	return lines, nil
}

func getExplorerDetailsDependingOnProduct(details PdDetails) (string, string) {
	// depending on service, return relevant explorer-information values and names for them
	// (contract address, key-hash, or registry address).
	// this function is also extra convoluted because alert details are not
	// standardized even within same-service alerts.
	switch details.CllService {
	case "feeds":
		if details.ContractAddress != "" {
			return "contract_address", details.ContractAddress
		}
		if details.Contract != "" {
			return "contract", details.Contract
		}
		return "", ""

	case "vrf":
		if details.KeyUnderscoreHash != "" {
			return "key_hash", details.KeyUnderscoreHash
		}
		if details.KeyHash != "" {
			return "keyhash", details.KeyHash
		}
		return "", ""

	case "keepers":
		return "registry_address", details.RegistryAddress

	default:
		return "", ""
	}
}

func IsUrl(str string) bool {
	// basic check to see if a string is a good URL
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}
