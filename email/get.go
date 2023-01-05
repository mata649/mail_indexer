package email

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func GetEmails(emailPaths []string, currentDir string) {
	counter := 0
	var emails []Email

	for _, path := range emailPaths {

		file, err := os.Open(path)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		var email Email
		currentParam := ""
		mainParamsChecked := false
		for scanner.Scan() {
			line := scanner.Text()
			firstLineParam := false
			var lineSplitted []string
			if !mainParamsChecked {

				lineSplitted = strings.Split(line, ":")
			}
			if len(lineSplitted) > 1 {
				currentParam = lineSplitted[0]
				if currentParam == "X-FileName" {

					mainParamsChecked = true
					continue
				}
				firstLineParam = true
			}
			if mainParamsChecked == false {
				switch currentParam {
				case "Message-ID":
					email.MessageID = strings.Trim(lineSplitted[1], " ")
				case "Date":
					email.Date = strings.Trim(strings.Join(lineSplitted[1:], ":"), " ")
				case "From":
					email.From = lineSplitted[1]
				case "To":
					if firstLineParam {
						emailsSplited := strings.Split(lineSplitted[1], ",")
						if len(emailsSplited) < 2 {
							email.To = append(email.To, strings.Trim(lineSplitted[1], " "))
							continue
						}
						for _, emailSplited := range emailsSplited {
							email.To = append(email.To, strings.Trim(emailSplited, " "))
						}
						continue
					}
					for _, emailSplited := range strings.Split(line, ",") {
						email.To = append(email.To, strings.Trim(strings.Trim(emailSplited, "	"), " "))

					}
				case "Subject":
					if firstLineParam {

						email.Subject = strings.Trim(strings.Join(lineSplitted[1:], ":"), " ")
						continue
					}
					email.Subject += lineSplitted[0]

				}
			} else if mainParamsChecked {
				email.Content += line + "\n"
			}

			if err := scanner.Err(); err != nil {
				fmt.Println(err)
			}
		}
		emails = append(emails, email)

		if len(emails) == 1000 {
			counter += 1
			saveEmails(emails, counter, currentDir)
			emails = []Email{}

		}
	}
	saveEmails(emails, counter, currentDir)

}
