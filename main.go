package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"reflect"
	"time"
)

const (
	EXIT_CODE     = 1
	STUCKDURATION = 0 //Hours
	VERSIONMSG    = "Veracode CLI version"
)
var Version = "source"

func main() {
	commands_desc := "[Mandatory] Veracode command: \n\t" +
		"getappid - Finding application ID by using application name. \n\t" +
		"getbuildinfo - Showing the information of the lastest build \n\t" +
		"getbuildversion - Retrieving the latest build version name \n\t" +
		"getbuildid - Retrieving the build id by using build name/version \n\t" +
		"buildstatus - Checking the status of the latest build and delete it if it's been stuck \n\t" +
		"devsecopscheck - Checking DevSecops check of the latest/specific build."
	userPtr := flag.String("user", "", "[Mandatory] Username")
	passPtr := flag.String("pass", "", "[Mandatory] Password")
	command := flag.String("command", "", commands_desc)
	appName := flag.String("appname", "", "[Mandatory] Application Name")
	buildName := flag.String("buildname", "", "[Optional] Build Name")
	versionPtr := flag.Bool("version", false, "Prints the current version.")
	flag.Parse()
	if *versionPtr {
		fmt.Fprintln(os.Stdout, VERSIONMSG, Version)
		os.Exit(0)
	}
	if *userPtr == "" || *passPtr == "" || *command == "" || *appName == "" {
		flag.PrintDefaults()
		os.Exit(EXIT_CODE)
	}

	var credentials VeracodeCredentials

	credentials.Username = *userPtr
	credentials.Password = *passPtr

	var arguments VeracodeArgs

	arguments.AppName = *appName
	if *buildName != "" {
		arguments.BuildName = *buildName
	}

	switch *command {
	case "getappid":
		if *appName != "" {
			log.Println(VERSIONMSG, Version)
			log.Println("getappid")
			err := FindAppIdByName(credentials, &arguments.AppName, &arguments.AppID)
			if err != nil {
				log.Println(err)
				os.Exit(EXIT_CODE)
			} else {
				fmt.Fprintln(os.Stdout, arguments.AppID)
			}
		} else {
			flag.PrintDefaults()
			os.Exit(1)
		}
	case "getbuildinfo":
		if *appName != "" {
			log.Println(VERSIONMSG, Version)
			log.Println("getbuildinfo")
			err := FindAppIdByName(credentials, &arguments.AppName, &arguments.AppID)
			if err != nil {
				log.Println(err)
				os.Exit(EXIT_CODE)
			} else {
				Binfo, err := getbuildinfo(credentials, arguments.AppID)
				if err != nil {
					log.Println(err)
					os.Exit(EXIT_CODE)
				}
				log.Println("Build ID: ", Binfo.Build.BuildID)
				log.Println("Build Version: ", Binfo.Build.Version)
				log.Println("Build Status: ", Binfo.Build.AnalysisUnit.Status)
			}
		} else {
			flag.PrintDefaults()
			os.Exit(1)
		}
	case "getbuildversion":
		if *appName != "" {
			log.Println(VERSIONMSG, Version)
			log.Println("getbuildversion")
			err := FindAppIdByName(credentials, &arguments.AppName, &arguments.AppID)
			if err != nil {
				log.Println(err)
				os.Exit(EXIT_CODE)
			} else {
				id, err := getbuildVersion(credentials, arguments.AppID)
				if err != nil {
					log.Println(err)
					os.Exit(EXIT_CODE)
				} else {
					fmt.Fprintln(os.Stdout, id)
				}
			}
		} else {
			flag.PrintDefaults()
			os.Exit(1)
		}

	case "getbuildid":
		if *appName != "" && *buildName != "" {
			log.Println(VERSIONMSG, Version)
			log.Println("getbuildid")
			err := FindAppIdByName(credentials, &arguments.AppName, &arguments.AppID)
			if err != nil {
				log.Println(err)
				os.Exit(EXIT_CODE)
			} else {
				err := FindBuildIdByBuildName(credentials, &arguments.BuildID, &arguments.AppID, &arguments.BuildName)
				if err != nil {
					log.Println(err)
					os.Exit(EXIT_CODE)
				} else {
					fmt.Fprintln(os.Stdout, arguments.BuildID)
				}
			}

		} else {
			flag.PrintDefaults()
			os.Exit(1)
		}

	case "buildstatus":
		if *appName != "" {
			log.Println(VERSIONMSG, Version)
			log.Println("Checking Build status")
			err := buildstatus(credentials, &arguments.AppName)
			if (err.Error() == SCAN_IS_READY) || (err.Error() == SCAN_STUCK_AND_DELETED) || (err.Error() == BUILD_NOT_FOUND){
				log.Println(err)
				fmt.Fprintln(os.Stdout, APP_IS_READY)
			} else {
				log.Println(err.Error())
				os.Exit(EXIT_CODE)
			}
		} else {
			flag.PrintDefaults()
			os.Exit(1)
		}

	case "devsecopscheck":
		if *appName != "" {
			log.Println(VERSIONMSG, Version)
			log.Println("DevSecopscheck")
			err := FindAppIdByName(credentials, &arguments.AppName, &arguments.AppID)
			if err != nil {
				if err.Error() == APP_NOT_FOUND {
					log.Println(err)
					fmt.Fprintln(os.Stdout, FLAG_APP_NOT_FOUND+"-"+err.Error())
				} else {
					log.Println(err)
					fmt.Fprintln(os.Stdout, FLAG_APP_ERROR+"-"+err.Error())
				}
				os.Exit(0)
			} else {
				// Checking specific build
				if *buildName != "" {
					err = FindBuildIdByBuildName(credentials, &arguments.BuildID, &arguments.AppID, &arguments.BuildName)
					if err != nil {
						if err.Error() == BUILD_NOT_FOUND {
							log.Println(err)
							fmt.Fprintln(os.Stdout, FLAG_BUILD_NOT_FOUND+"-"+err.Error())
						} else {
							log.Println(err)
							fmt.Fprintln(os.Stdout, FLAG_BUILD_ERROR+"-"+err.Error())
						}
						os.Exit(0)
					}
				}
				severities, status := DevSecopsCheck(credentials, arguments.AppID, arguments.BuildID)
				if !reflect.DeepEqual(severities, VeracodeSeverity{}) && status != nil {
					log.Println("[Result] { High & V.High:", severities.HighAndVeryHigh, ", Medium:", severities.Medium, "}")
					log.Println(status)
					fmt.Fprintln(os.Stdout, FLAG_APP_IS_NOT_OK+"-"+status.Error())
					os.Exit(0)
				}
				if status != nil {
					if status.Error() == STATUS_SCAN_IS_NOT_READY {
						log.Println(status)
						fmt.Fprintln(os.Stdout, FLAG_BUILD_NOT_READY+"-"+status.Error())
					} else {
						log.Println(status)
						fmt.Fprintln(os.Stdout, FLAG_REPORT_ERROR+"-"+STATUS_REPORT_UNAVAIL)
					}
					os.Exit(0)
				}
				log.Println(APP_IS_OK)
				fmt.Fprintln(os.Stdout, FLAG_APP_IS_OK+"-"+APP_IS_OK)
			}
		} else {
			flag.PrintDefaults()
			os.Exit(EXIT_CODE)
		}
	default:
		fmt.Println(INVALID_COMMAND)
		fmt.Println(commands_desc)
		os.Exit(EXIT_CODE)
	}
}

func getbuildinfo(credentials VeracodeCredentials, app_id string) (BuildInfo, error) {
	var Binfo BuildInfo
	err := VeracodeLastBuildInfo(credentials, &app_id, &Binfo)
	if err != nil {
		return Binfo, err
	}
	return Binfo, nil
}

func getbuildVersion(credentials VeracodeCredentials, app_id string) (string, error) {
	var Binfo BuildInfo
	err := VeracodeLastBuildInfo(credentials, &app_id, &Binfo)
	if err != nil {
		return "", err
	}
	return Binfo.Build.Version, nil
}

func DevSecopsCheck(credentials VeracodeCredentials, app_id string, build_id string) (VeracodeSeverity, error) {
	var Binfo BuildInfo
	var severity_total VeracodeSeverity
	var err error
	if build_id != "" {
		err = VeracodeBuildInfo(credentials, &app_id, &build_id, &Binfo)
	} else {
		err = VeracodeLastBuildInfo(credentials, &app_id, &Binfo)
	}
	if err != nil {
		return severity_total, err
	}
	//Checking the build status
	err = ScanCheckStatus(&Binfo)
	if err != nil {
		//Analyze the report if the build is ready
		//APP_IS_NOT_OK means the report is exist but has some severities
		if err.Error() == APP_IS_NOT_OK {
			// Download Full report
			report, err := downloadFullReport(credentials, &Binfo.Build.BuildID)
			if err != nil {
				return severity_total, err
			}
			severity_total, _ = SeveritiesNotApproved(&report)
			if severity_total.Medium+severity_total.HighAndVeryHigh != 0 {
				return severity_total, errors.New(APP_IS_NOT_OK)
			}
		} else {
			return severity_total, errors.New(STATUS_SCAN_IS_NOT_READY)
		}
	}
	return severity_total, nil
}

func getAppID(credentials VeracodeCredentials, app_name *string, app_id *string) error {
	err := FindAppIdByName(credentials, app_name, app_id)
	if err != nil {
		return errors.New(APP_NOT_FOUND)
	}
	return err
}

func buildstatus(credentials VeracodeCredentials, app_name *string) error {
	var Binfo BuildInfo
	var app_id string
	err := FindAppIdByName(credentials, app_name, &app_id)
	if err != nil {
		return err
	}

	err = VeracodeLastBuildInfo(credentials, &app_id, &Binfo)
	if err != nil {
		return err
	}
	err = ScanCheckStatus(&Binfo)

	if err == nil || err.Error() == APP_IS_NOT_OK {
		return errors.New(SCAN_IS_READY)
	}

	if err.Error() == STATUS_SCAN_INCOMPLETE || err.Error() == STATUS_PRE_SCAN_FAILED || 
	err.Error() == STATUS_PRE_SCAN_SUBMITTED || err.Error() == STATUS_NO_MODULES_DEFINED {
		err = deleteAppLastBuild(credentials, app_id)
		if err != nil {
			return err
		}
		return errors.New(SCAN_STUCK_AND_DELETED)
	}

	if err.Error() == SCAN_IS_IN_PROGRESS || err.Error() == STATUS_PRE_SCAN_SUCCESS{
		if Binfo.Build.PolicyUpdatedDate != "" {
			elapsed, err := TimeSinceLastScan(&Binfo.Build.PolicyUpdatedDate)
			if err != nil {
				return err
			}
			if elapsed >= STUCKDURATION {
				err = deleteAppLastBuild(credentials, app_id)
				if err != nil {
					return err
				}
				return errors.New(SCAN_STUCK_AND_DELETED)
			} else {
				return errors.New(SCAN_STUCKEDNOTDELETED)
			}
		} else {
			return errors.New(SCAN_STUCKED)
		}
	}
	return err
}

func TimeSinceLastScan(scantime *string) (int, error) {
	t2, err := time.Parse(time.RFC3339, *scantime)
	LastScan := int(time.Since(t2).Hours())
	if err != nil {
		return -1, err
	}
	return LastScan, nil
}
