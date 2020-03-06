package pfsagentConfig

import (
	"fmt"
	"log"
	"os"
)

func RunStateMachine() error {
	// fmt.Println("runWizard starting")
	nextMenuText := mainMenuText
	nextMenuOptions := mainMenuOptions
	nextMenuOptionsTexts := mainMenuOptionsTexts
	prevMenuText := mainMenuText
	prevMenuOptions := mainMenuOptions
	prevMenuOptionsTexts := mainMenuOptionsTexts
	for true {
		menuResponse, displayErr := nextMenu(nextMenuText, nextMenuOptions, nextMenuOptionsTexts)
		// fmt.Printf("menuResponse: %v\n", menuResponse)
		if nil != displayErr {
			fmt.Println("ERROR while displaying menu item", displayErr)
			return fmt.Errorf("error trying to display menu item")
		}
		switch menuResponse {
		case quitMenuOptionText:
			fmt.Println("Thank you for using the pfsagent config util")
			return nil
		case backMenuOptionText:
			nextMenuText = prevMenuText
			nextMenuOptions = prevMenuOptions
			nextMenuOptionsTexts = prevMenuOptionsTexts
		case changeCredsOptionText:
			// fmt.Printf("got %v\n", changeCredsOptionText)
			prevMenuText = nextMenuText
			prevMenuOptions = nextMenuOptions
			prevMenuOptionsTexts = nextMenuOptionsTexts
			nextMenuText = credentialsMenuTexts
			nextMenuOptions = credentialsMenuOptions
			nextMenuOptionsTexts = credentialsMenuOptionsTexts
		case changeOtherOptionText:
			fmt.Printf("got %v\n", changeOtherOptionText)

		case changeAuthURLOptionText:
			userResponse, userInputErr := getValueFromUser("Swift Auth URL", "", confMap["Agent"]["SwiftAuthURL"][0])
			if nil != userInputErr {
				fmt.Println("Error reading input from user", userInputErr)
				return userInputErr
			}
			prevAuthURL := confMap["Agent"]["SwiftAuthURL"][0]
			confMap["Agent"]["SwiftAuthURL"][0] = userResponse
			whatFailed, accessErr := ValidateAccess()
			if nil != accessErr {
				switch whatFailed {
				case typeAuthURL:
					confMap["Agent"]["SwiftAuthURL"][0] = prevAuthURL
					fmt.Printf(failureMessageHeader)
					fmt.Printf(authURLFailedMessage, accessErr)
					fmt.Printf(failureMessageFooter)
					continue
				case typeCredentails:
					fmt.Printf(needMoreInfoMessageHeader)
					fmt.Printf("Auth URL Works, But I Got An Error Trying To Login With Credentails\nUser: %v\nKey: %v\n%v\n\n", confMap["Agent"]["SwiftAuthUser"][0], confMap["Agent"]["SwiftAuthKey"][0], accessErr)
					// MyConfig.SwiftAuthURL = userResponse
					fmt.Printf("Swift Auth URL Set To %v\n", confMap["Agent"]["SwiftAuthURL"][0])
					fmt.Printf(needMoreInfoMessageFooter)
					continue
				case typeAccount:
					fmt.Printf(needMoreInfoMessageHeader)
					fmt.Printf("Auth URL And Credentials Works, But I Could Not Gain Access To Account %v. Please Verify The Account Exists And User %v Has The Correct Access Permissions\n%v\n\n", confMap["Agent"]["SwiftAccountName"][0], confMap["Agent"]["SwiftAuthUser"][0], accessErr)
					// MyConfig.SwiftAuthURL = userResponse
					fmt.Printf("Swift Auth URL Set To %v\n", confMap["Agent"]["SwiftAuthURL"][0])
					SaveCurrentConfig()
					fmt.Println("Changes Saved To File")
					fmt.Printf(needMoreInfoMessageFooter)
					continue
				}
			} else {
				fmt.Printf(successMessageHeader)
				fmt.Printf("All Access Checks Succeeded")
				// MyConfig.SwiftAuthURL = userResponse
				fmt.Printf("Swift Auth URL Set To %v\n", confMap["Agent"]["SwiftAuthURL"][0])
				SaveCurrentConfig()
				fmt.Println("Changes Saved To File")
				fmt.Printf(successMessageFooter)
				nextMenuText = mainMenuText
				nextMenuOptions = mainMenuOptions
				nextMenuOptionsTexts = mainMenuOptionsTexts
				continue
			}

		case changeUsernameOptionText:
			userResponse, userInputErr := getValueFromUser("Swift Username", "", confMap["Agent"]["SwiftAuthUser"][0])
			if nil != userInputErr {
				fmt.Println("Error reading input from user", userInputErr)
				return userInputErr
			}
			prevAuthUser := confMap["Agent"]["SwiftAuthUser"][0]
			confMap["Agent"]["SwiftAuthUser"][0] = userResponse
			whatFailed, accessErr := ValidateAccess()
			if nil != accessErr {
				switch whatFailed {
				case typeAuthURL:
					confMap["Agent"]["SwiftAuthUser"][0] = prevAuthUser
					fmt.Printf(failureMessageHeader)
					fmt.Printf("Auth URL Failed, So I Could Not Check Username. Please Verify Auth URL\n%v\n\n", accessErr)
					fmt.Printf(failureMessageFooter)
					continue
				case typeCredentails:
					// MyConfig.SwiftAuthUser = prevAuthUser
					fmt.Printf("Swift User Set To %v\n", confMap["Agent"]["SwiftAuthUser"][0])
					SaveCurrentConfig()
					fmt.Println("Changes Saved To File")
					fmt.Printf(failureMessageHeader)
					fmt.Printf("Auth URL Works, But I Got An Error Trying To Login With Credentails\nUser: %v\nKey: %v\n%v\n\n", userResponse, confMap["Agent"]["SwiftAuthKey"][0], accessErr)
					fmt.Printf(failureMessageFooter)
					continue
				case typeAccount:
					fmt.Printf(needMoreInfoMessageHeader)
					fmt.Printf("Auth URL And Credentials Works, But I Could Not Gain Access To Account %v. Please Verify The Account Exists And User %v Has The Correct Access Permissions\n%v\n\n", confMap["Agent"]["SwiftAccountName"][0], confMap["Agent"]["SwiftAuthUser"][0], accessErr)
					// MyConfig.SwiftAuthUser = userResponse
					fmt.Printf("Swift User Set To %v\n", confMap["Agent"]["SwiftAuthUser"][0])
					SaveCurrentConfig()
					fmt.Println("Changes Saved To File")
					fmt.Printf(needMoreInfoMessageFooter)
					continue
				}
			} else {
				fmt.Printf(successMessageHeader)
				fmt.Printf("All Access Checks Succeeded\n")
				// MyConfig.SwiftAuthUser = userResponse
				fmt.Printf("Swift User Set To %v\n", confMap["Agent"]["SwiftAuthUser"][0])
				SaveCurrentConfig()
				fmt.Println("Changes Saved To File")
				fmt.Printf(successMessageFooter)
				nextMenuText = mainMenuText
				nextMenuOptions = mainMenuOptions
				nextMenuOptionsTexts = mainMenuOptionsTexts
				continue
			}

		case changeKeyOptionText:
			userResponse, userInputErr := getValueFromUser("Swift User Key", "", confMap["Agent"]["SwiftAuthKey"][0])
			if nil != userInputErr {
				fmt.Println("Error reading input from user", userInputErr)
				return userInputErr
			}
			prevAuthKey := confMap["Agent"]["SwiftAuthKey"][0]
			confMap["Agent"]["SwiftAuthKey"][0] = userResponse
			whatFailed, accessErr := ValidateAccess()
			if nil != accessErr {
				switch whatFailed {
				case typeAuthURL:
					confMap["Agent"]["SwiftAuthKey"][0] = prevAuthKey
					fmt.Printf(failureMessageHeader)
					fmt.Printf("Auth URL Failed, So I Could Not Check User Key. Please Verify Auth URL\n%v\n\n", accessErr)
					fmt.Printf(failureMessageFooter)
					continue
				case typeCredentails:
					// MyConfig.SwiftAuthKey = prevAuthKey
					fmt.Printf("Swift User Key To %v\n", confMap["Agent"]["SwiftAuthKey"][0])
					SaveCurrentConfig()
					fmt.Println("Changes Saved To File")
					fmt.Printf(failureMessageHeader)
					fmt.Printf("Auth URL Works, But I Got An Error Trying To Login With Credentails\nUser: %v\nKey: %v\n%v\n\n", confMap["Agent"]["SwiftAuthUser"][0], userResponse, accessErr)
					fmt.Printf(failureMessageFooter)
					continue
				case typeAccount:
					fmt.Printf(needMoreInfoMessageHeader)
					fmt.Printf("Auth URL And Credentials Works, But I Could Not Gain Access To Account %v. Please Verify The Account Exists And User %v Has The Correct Access Permissions\n%v\n\n", confMap["Agent"]["SwiftAccountName"][0], confMap["Agent"]["SwiftAuthUser"][0], accessErr)
					// MyConfig.SwiftAuthKey = userResponse
					fmt.Printf("Swift User Key Set to %v\n", confMap["Agent"]["SwiftAuthKey"][0])
					SaveCurrentConfig()
					fmt.Println("Changes Saved To File")
					fmt.Printf(needMoreInfoMessageFooter)
					continue
				}
			} else {
				fmt.Printf(successMessageHeader)
				fmt.Printf("All Access Checks Succeeded\n")
				// MyConfig.SwiftAuthKey = userResponse
				fmt.Printf("Swift User Key to %v\n", confMap["Agent"]["SwiftAuthKey"][0])
				SaveCurrentConfig()
				fmt.Println("Changes Saved To File")
				fmt.Printf(successMessageFooter)
				nextMenuText = mainMenuText
				nextMenuOptions = mainMenuOptions
				nextMenuOptionsTexts = mainMenuOptionsTexts
				continue
			}

		case changeAccountOptionText:
			userResponse, userInputErr := getValueFromUser("Swift Account", "", confMap["Agent"]["SwiftAccountName"][0])
			if nil != userInputErr {
				fmt.Println("Error reading input from user", userInputErr)
				return userInputErr
			}
			prevAccountName := confMap["Agent"]["SwiftAccountName"][0]
			confMap["Agent"]["SwiftAccountName"][0] = userResponse
			whatFailed, accessErr := ValidateAccess()
			if nil != accessErr {
				switch whatFailed {
				case typeAuthURL:
					confMap["Agent"]["SwiftAccountName"][0] = prevAccountName
					fmt.Printf(failureMessageHeader)
					fmt.Printf("Auth URL Failed, So I Could Not Check User Key. Please Verify Auth URL\n%v\n\n", accessErr)
					fmt.Printf(failureMessageFooter)
					continue
				case typeCredentails:
					confMap["Agent"]["SwiftAccountName"][0] = prevAccountName
					fmt.Printf(failureMessageHeader)
					fmt.Printf("Auth URL Works, But I Got An Error Trying To Login With Credentails\nUser: %v\nKey: %v\n%v\n\n", confMap["Agent"]["SwiftAuthUser"][0], confMap["Agent"]["SwiftAuthKey"][0], accessErr)
					fmt.Printf(failureMessageFooter)
					continue
				case typeAccount:
					confMap["Agent"]["SwiftAccountName"][0] = prevAccountName
					fmt.Printf(failureMessageHeader)
					fmt.Printf("Auth URL And Credentials Works, But I Could Not Gain Access To Account %v. Please Verify The Account Exists And User %v Has The Correct Access Permissions\n%v\n\n", confMap["Agent"]["SwiftAccountName"][0], confMap["Agent"]["SwiftAuthUser"][0], accessErr)
					fmt.Printf(failureMessageFooter)
					continue
				}
			} else {
				fmt.Printf(successMessageHeader)
				fmt.Printf("All Access Checks Succeeded\n")
				// MyConfig.SwiftAccountName = userResponse
				fmt.Printf("Swift Account Set To %v\n", confMap["Agent"]["SwiftAccountName"][0])
				SaveCurrentConfig()
				fmt.Println("Changes Saved To File")
				fmt.Printf(successMessageFooter)
				nextMenuText = mainMenuText
				nextMenuOptions = mainMenuOptions
				nextMenuOptionsTexts = mainMenuOptionsTexts
				continue
			}

		default:
			fmt.Printf("got unknown response: %v\n", menuResponse)
			continue
		}
	}
	fmt.Println("State machine exiting")
	return nil
}

func FirstTimeRun() error {
	loadError := LoadConfig("")
	if nil != loadError {
		fmt.Println("Failed loading config. Error:", loadError)
		os.Exit(1)
	}

	var oldAuthURL string
	var oldAuthUser string
	var oldAuthKey string
	var oldAccount string
	var oldMount string
	var oldVolName string
	var oldLogPath string

	if len(confMap["Agent"]["SwiftAuthURL"]) > 0 {
		oldAuthURL = confMap["Agent"]["SwiftAuthURL"][0]
		log.Printf("past %v\n", oldAuthURL)
	}
	if len(confMap["Agent"]["SwiftAuthUser"]) > 0 {
		oldAuthUser = confMap["Agent"]["SwiftAuthUser"][0]
		log.Printf("past %v\n", oldAuthUser)
	}
	if len(confMap["Agent"]["SwiftAuthKey"]) > 0 {
		oldAuthKey = confMap["Agent"]["SwiftAuthKey"][0]
		log.Printf("past %v\n", oldAuthKey)
	}
	if len(confMap["Agent"]["SwiftAccountName"]) > 0 {
		oldAccount = confMap["Agent"]["SwiftAccountName"][0]
		log.Printf("past %v\n", oldAccount)
	}
	if len(confMap["Agent"]["FUSEMountPointPath"]) > 0 {
		oldMount = confMap["Agent"]["FUSEMountPointPath"][0]
		log.Printf("past %v\n", oldMount)
	}
	if len(confMap["Agent"]["FUSEVolumeName"]) > 0 {
		oldVolName = confMap["Agent"]["FUSEVolumeName"][0]
		log.Printf("past %v\n", oldVolName)
	}
	if len(confMap["Agent"]["LogFilePath"]) > 0 {
		oldLogPath = confMap["Agent"]["LogFilePath"][0]
		log.Printf("past %v\n", oldLogPath)
	}

	confMap["Agent"]["LogFilePath"][0] = "ff"

	SaveCurrentConfig()

	fmt.Printf("conf map:\n%v\n", confMap)

	fmt.Println(firstTimeCredentialsMenu)
	userURLResponse, userURLInputErr := getValueFromUser("Swift Auth URL", "", "")
	if nil != userURLInputErr {
		return fmt.Errorf("Error Reading Auth URL From User\n%v", userURLInputErr)
	}

	userUserResponse, userUserInputErr := getValueFromUser("Swift Auth User", "", "")
	if nil != userUserInputErr {
		return fmt.Errorf("Error Reading Auth User From User\n%v", userUserInputErr)
	}

	userKeyResponse, userKeyInputErr := getValueFromUser("Swift Auth Key", "", "")
	if nil != userKeyInputErr {
		return fmt.Errorf("Error Reading Auth Key From User\n%v", userKeyInputErr)
	}

	userAccountResponse, userAccountInputErr := getValueFromUser("Swift Account", "", "")
	if nil != userAccountInputErr {
		return fmt.Errorf("Error Reading Swift Account From User\n%v", userAccountInputErr)
	}

	volNameResponse, volNameInputErr := getValueFromUser("Volume Name", "", userAccountResponse)
	if nil != volNameInputErr {
		return fmt.Errorf("Error Reading Volume Name From User\n%v", volNameInputErr)
	}

	var suggestedMountPath = fmt.Sprintf("%v/%v", defaultMountPath, userAccountResponse)
	var suggestedLogPath = fmt.Sprintf("%v/%v", defaultLogPath, userAccountResponse)

	mountPathResponse, mountPathInputErr := getValueFromUser("Mount Point", "", suggestedMountPath)
	if nil != mountPathInputErr {
		return fmt.Errorf("Error Reading Mount Path From User\n%v", mountPathInputErr)
	}

	confMap["Agent"]["SwiftAuthURL"][0] = userURLResponse
	confMap["Agent"]["SwiftAuthUser"][0] = userUserResponse
	confMap["Agent"]["SwiftAuthKey"][0] = userKeyResponse
	confMap["Agent"]["SwiftAccountName"][0] = userAccountResponse
	if len(mountPathResponse) > 0 {
		confMap["Agent"]["FUSEMountPointPath"][0] = mountPathResponse
	} else {
		confMap["Agent"]["FUSEMountPointPath"][0] = suggestedMountPath
	}
	if len(volNameResponse) > 0 {
		confMap["Agent"]["FUSEVolumeName"][0] = userAccountResponse
	} else {
		confMap["Agent"]["FUSEVolumeName"][0] = volNameResponse
	}

	confMap["Agent"]["LogFilePath"][0] = suggestedLogPath

	whatFailed, accessErr := ValidateAccess()
	if nil != accessErr {
		confMap["Agent"]["SwiftAuthURL"][0] = oldAuthURL
		confMap["Agent"]["SwiftAuthUser"][0] = oldAuthUser
		confMap["Agent"]["SwiftAuthKey"][0] = oldAuthKey
		confMap["Agent"]["SwiftAccountName"][0] = oldAccount
		confMap["Agent"]["FUSEMountPointPath"][0] = oldMount
		confMap["Agent"]["FUSEVolumeName"][0] = oldVolName
		confMap["Agent"]["LogFilePath"][0] = oldLogPath
		switch whatFailed {
		case typeAuthURL:
			fmt.Printf(failureMessageHeader)
			fmt.Printf(authURLFailedMessage, accessErr)
			fmt.Printf(failureMessageFooter)
		case typeCredentails:
			fmt.Printf(failureMessageHeader)
			fmt.Printf("Auth URL Works, But I Got An Error Trying To Login With Credentails\nUser: %v\nKey: %v\n%v\n\n", confMap["Agent"]["SwiftAuthUser"][0], confMap["Agent"]["SwiftAuthKey"][0], accessErr)
			fmt.Printf("Swift Auth URL Set To %v\n", confMap["Agent"]["SwiftAuthURL"][0])
			fmt.Printf(failureMessageFooter)
		case typeAccount:
			fmt.Printf(failureMessageHeader)
			fmt.Printf("Auth URL And Credentials Works, But I Could Not Gain Access To Account %v. Please Verify The Account Exists And User %v Has The Correct Access Permissions\n%v\n\n", confMap["Agent"]["SwiftAccountName"][0], confMap["Agent"]["SwiftAuthUser"][0], accessErr)
			fmt.Printf("Swift Auth URL Set To %v\n", confMap["Agent"]["SwiftAuthURL"][0])
			SaveCurrentConfig()
			fmt.Println("Changes Saved To File")
			fmt.Printf(failureMessageFooter)
		}
	} else {
		fmt.Printf(successMessageHeader)
		fmt.Printf("All Access Checks Succeeded")

		if _, err := os.Stat(confMap["Agent"]["LogFilePath"][0]); os.IsNotExist(err) {
			err = os.MkdirAll(confMap["Agent"]["LogFilePath"][0], 0755)
			if err != nil {
				fmt.Printf(failureMessageHeader)
				// fmt.Printf(authURLFailedMessage, accessErr)
				fmt.Printf(failureMessageFooter)
				panic(err)
			}
		}

		if _, err := os.Stat(confMap["Agent"]["FUSEMountPointPath"][0]); os.IsNotExist(err) {
			err = os.MkdirAll(confMap["Agent"]["FUSEMountPointPath"][0], 0755)
			if err != nil {
				fmt.Printf(failureMessageHeader)
				// fmt.Printf(authURLFailedMessage, accessErr)
				fmt.Printf(failureMessageFooter)
				panic(err)
			}
		}
		SaveCurrentConfig()
		fmt.Println("Changes Saved To File")
		fmt.Printf(successMessageFooter)
	}

	return nil
}
