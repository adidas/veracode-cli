package main

import (
	"encoding/xml"
	"errors"
	"log"
)

const (
	APP_IS_OK         = "[success] The build has no none-mitigated flaws with high severity."
	APP_IS_NOT_OK     = "[Error] The build has some none-mitigated flaws with high severity!"
 	APP_IS_READY              = "[Ok] App is ready for new scan"
	SCAN_IS_READY             = "[Ok] Results Ready"
	SCAN_IS_IN_PROGRESS       = "[Error] The scan is still in progress. Please wait !"
	SCAN_STUCKED              = "[Error] The scan has been stuck and cannot be deleted automatically !"
	SCAN_STUCK_AND_DELETED    = "[Success] The scan was stuck and successfully deleted."
	SCAN_STUCKEDNOTDELETED    = "[Error] The scan has been stuck less than the threshold! Please try again later!"
	APP_NOT_FOUND             = "[Error] The App was not found!"
	BUILD_NOT_FOUND           = "[Error] The Build was not found!"
	CATEGORY_IS_EMPTY         = "[Error] Category is empty"
	API_CON_ERROR             = "[Error] Veracode API connection issue!"
	API_UNEXPECTED_ERROR	  = "[Error] Unexpected error on API response!"
	API_NOT_FOUND             = "[Error] The requested API was not found!"
	REQ_BUILDLIST             = "Requesting build list"
	REQ_BUILD_INFO            = "Requesting build info"
	REQ_LAST_BUILD            = "Requesting last build"
	REQ_DELETE_BUILD          = "Deleting last build"
	REQ_FULL_REPORT           = "Requesting full report"
	REQ_FIND_APPID            = "Finding App ID"
	BUILD_STATUS              = "The build status is [%s]"
	STATUS_SCAN_IS_READY      = "Results Ready"
	STATUS_SCAN_IS_NOT_READY  = "[Error] Scan result is not ready!"
	STATUS_SCAN_IN_PROGRESS   = "Scan In Process"
	STATUS_SCAN_INCOMPLETE    = "Incomplete"
	STATUS_PRE_SCAN_FAILED    = "Pre-Scan Failed"
	STATUS_PRE_SCAN_SUCCESS   = "Pre-Scan Success"
	STATUS_PRE_SCAN_SUBMITTED = "Pre-Scan Submitted"
	STATUS_NO_MODULES_DEFINED = "No Modules Defined"
	STATUS_SCAN_UNKNOWN       = "[Error] Unknown status: "
	STATUS_REPORT_UNAVAIL     = "No report available."
	FLAG_APP_NOT_FOUND		 = "N000"
	FLAG_APP_ERROR		  	 = "0000"
	FLAG_BUILD_NOT_FOUND	 = "YN00"
	FLAG_BUILD_ERROR		   = "Y000"
	FLAG_APP_IS_NOT_OK     = "YYYY"
	FLAG_APP_IS_OK		     = "YYYN"
	FLAG_BUILD_NOT_READY	 = "YYN0"
	FLAG_REPORT_ERROR		   = "YYY0"
	INVALID_COMMAND        = "[Error] Invalid Command."
	INVALID_CREDS          = "[Error] Invalid Credentials. Authentication failed !"
	FORBIDEN_ACCESS        = "[Error] Your account doesn't have sufficient access. Please contact with Security team."
)

type VeracodeApplist struct {
	XMLName        xml.Name `xml:"applist"`
	Text           string   `xml:",chardata"`
	Xsi            string   `xml:"xsi,attr"`
	Xmlns          string   `xml:"xmlns,attr"`
	SchemaLocation string   `xml:"schemaLocation,attr"`
	ApplistVersion string   `xml:"applist_version,attr"`
	AccountID      string   `xml:"account_id,attr"`
	App            []struct {
		Text              string `xml:",chardata"`
		AppID             string `xml:"app_id,attr"`
		AppName           string `xml:"app_name,attr"`
		PolicyUpdatedDate string `xml:"policy_updated_date,attr"`
	} `xml:"app"`
}

type VeracodeBuildList struct {
	Buildlist        xml.Name `xml:"buildlist"`
	Text             string   `xml:",chardata"`
	Xsi              string   `xml:"xsi,attr"`
	Xmlns            string   `xml:"xmlns,attr"`
	SchemaLocation   string   `xml:"schemaLocation,attr"`
	BuildlistVersion string   `xml:"buildlist_version,attr"`
	AccountID        string   `xml:"account_id,attr"`
	AppID            string   `xml:"app_id,attr"`
	AppName          string   `xml:"app_name,attr"`
	Build            []struct {
		Text              string `xml:",chardata"`
		BuildID           string `xml:"build_id,attr"`
		Version           string `xml:"version,attr"`
		PolicyUpdatedDate string `xml:"policy_updated_date,attr"`
		DynamicScanType   string `xml:"dynamic_scan_type,omitempty,attr"`
	} `xml:"build"`
}

type BuildInfo struct {
	XMLName          xml.Name `xml:"buildinfo"`
	Text             string   `xml:",chardata"`
	Xsi              string   `xml:"xsi,attr"`
	Xmlns            string   `xml:"xmlns,attr"`
	SchemaLocation   string   `xml:"schemaLocation,attr"`
	BuildinfoVersion string   `xml:"buildinfo_version,attr"`
	AccountID        string   `xml:"account_id,attr"`
	AppID            string   `xml:"app_id,attr"`
	BuildID          string   `xml:"build_id,attr"`
	Build            struct {
		Text                   string `xml:",chardata"`
		Version                string `xml:"version,attr"`
		BuildID                string `xml:"build_id,attr"`
		Submitter              string `xml:"submitter,attr"`
		Platform               string `xml:"platform,attr"`
		LifecycleStage         string `xml:"lifecycle_stage,attr"`
		ResultsReady           string `xml:"results_ready,attr"`
		PolicyName             string `xml:"policy_name,attr"`
		PolicyVersion          string `xml:"policy_version,attr"`
		PolicyComplianceStatus string `xml:"policy_compliance_status,attr"`
		PolicyUpdatedDate      string `xml:"policy_updated_date,attr"`
		RulesStatus            string `xml:"rules_status,attr"`
		GracePeriodExpired     string `xml:"grace_period_expired,attr"`
		ScanOverdue            string `xml:"scan_overdue,attr"`
		LegacyScanEngine       string `xml:"legacy_scan_engine,attr"`
		AnalysisUnit           struct {
			Text             string `xml:",chardata"`
			AnalysisType     string `xml:"analysis_type,attr"`
			PublishedDate    string `xml:"published_date,attr"`
			PublishedDateSec string `xml:"published_date_sec,attr"`
			Status           string `xml:"status,attr"`
			EngineVersion    string `xml:"engine_version,attr"`
		} `xml:"analysis_unit"`
	} `xml:"build"`
}

type Detailedreport struct {
	XMLName                 xml.Name `xml:"detailedreport"`
	Text                    string   `xml:",chardata"`
	Xsi                     string   `xml:"xsi,attr"`
	Xmlns                   string   `xml:"xmlns,attr"`
	SchemaLocation          string   `xml:"schemaLocation,attr"`
	ReportFormatVersion     string   `xml:"report_format_version,attr"`
	AccountID               string   `xml:"account_id,attr"`
	AppName                 string   `xml:"app_name,attr"`
	AppID                   string   `xml:"app_id,attr"`
	AnalysisID              string   `xml:"analysis_id,attr"`
	StaticAnalysisUnitID    string   `xml:"static_analysis_unit_id,attr"`
	SandboxID               string   `xml:"sandbox_id,attr"`
	FirstBuildSubmittedDate string   `xml:"first_build_submitted_date,attr"`
	Version                 string   `xml:"version,attr"`
	BuildID                 string   `xml:"build_id,attr"`
	Submitter               string   `xml:"submitter,attr"`
	Platform                string   `xml:"platform,attr"`
	AssuranceLevel          string   `xml:"assurance_level,attr"`
	BusinessCriticality     string   `xml:"business_criticality,attr"`
	GenerationDate          string   `xml:"generation_date,attr"`
	VeracodeLevel           string   `xml:"veracode_level,attr"`
	TotalFlaws              string   `xml:"total_flaws,attr"`
	FlawsNotMitigated       string   `xml:"flaws_not_mitigated,attr"`
	Teams                   string   `xml:"teams,attr"`
	LifeCycleStage          string   `xml:"life_cycle_stage,attr"`
	PlannedDeploymentDate   string   `xml:"planned_deployment_date,attr"`
	LastUpdateTime          string   `xml:"last_update_time,attr"`
	IsLatestBuild           string   `xml:"is_latest_build,attr"`
	PolicyName              string   `xml:"policy_name,attr"`
	PolicyVersion           string   `xml:"policy_version,attr"`
	PolicyComplianceStatus  string   `xml:"policy_compliance_status,attr"`
	PolicyRulesStatus       string   `xml:"policy_rules_status,attr"`
	GracePeriodExpired      string   `xml:"grace_period_expired,attr"`
	ScanOverdue             string   `xml:"scan_overdue,attr"`
	BusinessOwner           string   `xml:"business_owner,attr"`
	BusinessUnit            string   `xml:"business_unit,attr"`
	Tags                    string   `xml:"tags,attr"`
	LegacyScanEngine        string   `xml:"legacy_scan_engine,attr"`
	StaticAnalysis          struct {
		Text              string `xml:",chardata"`
		Rating            string `xml:"rating,attr"`
		Score             string `xml:"score,attr"`
		SubmittedDate     string `xml:"submitted_date,attr"`
		PublishedDate     string `xml:"published_date,attr"`
		Version           string `xml:"version,attr"`
		NextScanDue       string `xml:"next_scan_due,attr"`
		AnalysisSizeBytes string `xml:"analysis_size_bytes,attr"`
		EngineVersion     string `xml:"engine_version,attr"`
		Modules           struct {
			Text   string `xml:",chardata"`
			Module struct {
				Text         string `xml:",chardata"`
				Name         string `xml:"name,attr"`
				Compiler     string `xml:"compiler,attr"`
				Os           string `xml:"os,attr"`
				Architecture string `xml:"architecture,attr"`
				Loc          string `xml:"loc,attr"`
				Score        string `xml:"score,attr"`
				Numflawssev0 string `xml:"numflawssev0,attr"`
				Numflawssev1 string `xml:"numflawssev1,attr"`
				Numflawssev2 string `xml:"numflawssev2,attr"`
				Numflawssev3 string `xml:"numflawssev3,attr"`
				Numflawssev4 string `xml:"numflawssev4,attr"`
				Numflawssev5 string `xml:"numflawssev5,attr"`
			} `xml:"module"`
		} `xml:"modules"`
	} `xml:"static-analysis"`
	Severity []struct {
		Text     string     `xml:",chardata"`
		Level    string     `xml:"level,attr"`
		Category []Category `xml:"category"`
	} `xml:"severity"`
	FlawStatus struct {
		Text            string `xml:",chardata"`
		New             string `xml:"new,attr"`
		Reopen          string `xml:"reopen,attr"`
		Open            string `xml:"open,attr"`
		CannotReproduce string `xml:"cannot-reproduce,attr"`
		Fixed           string `xml:"fixed,attr"`
		Total           string `xml:"total,attr"`
		NotMitigated    string `xml:"not_mitigated,attr"`
		Sev1Change      string `xml:"sev-1-change,attr"`
		Sev2Change      string `xml:"sev-2-change,attr"`
		Sev3Change      string `xml:"sev-3-change,attr"`
		Sev4Change      string `xml:"sev-4-change,attr"`
		Sev5Change      string `xml:"sev-5-change,attr"`
	} `xml:"flaw-status"`
	Customfields struct {
		Text        string `xml:",chardata"`
		Customfield []struct {
			Text  string `xml:",chardata"`
			Name  string `xml:"name,attr"`
			Value string `xml:"value,attr"`
		} `xml:"customfield"`
	} `xml:"customfields"`
	SoftwareCompositionAnalysis struct {
		Text                     string `xml:",chardata"`
		ThirdPartyComponents     string `xml:"third_party_components,attr"`
		ViolatePolicy            string `xml:"violate_policy,attr"`
		ComponentsViolatedPolicy string `xml:"components_violated_policy,attr"`
		VulnerableComponents     struct {
			Text      string `xml:",chardata"`
			Component []struct {
				Text                             string `xml:",chardata"`
				ComponentID                      string `xml:"component_id,attr"`
				FileName                         string `xml:"file_name,attr"`
				Sha1                             string `xml:"sha1,attr"`
				AttrVulnerabilities              string `xml:"vulnerabilities,attr"`
				MaxCvssScore                     string `xml:"max_cvss_score,attr"`
				Version                          string `xml:"version,attr"`
				Library                          string `xml:"library,attr"`
				Vendor                           string `xml:"vendor,attr"`
				Description                      string `xml:"description,attr"`
				AddedDate                        string `xml:"added_date,attr"`
				ComponentAffectsPolicyCompliance string `xml:"component_affects_policy_compliance,attr"`
				New                              string `xml:"new,attr"`
				FilePaths                        struct {
					Text     string `xml:",chardata"`
					FilePath struct {
						Text  string `xml:",chardata"`
						Value string `xml:"value,attr"`
					} `xml:"file_path"`
				} `xml:"file_paths"`
				Licenses struct {
					Text    string `xml:",chardata"`
					License []struct {
						Text       string `xml:",chardata"`
						Name       string `xml:"name,attr"`
						SpdxID     string `xml:"spdx_id,attr"`
						LicenseURL string `xml:"license_url,attr"`
						RiskRating string `xml:"risk_rating,attr"`
					} `xml:"license"`
				} `xml:"licenses"`
				Vulnerabilities struct {
					Text          string `xml:",chardata"`
					Vulnerability []struct {
						Text                                 string `xml:",chardata"`
						CveID                                string `xml:"cve_id,attr"`
						CvssScore                            string `xml:"cvss_score,attr"`
						Severity                             string `xml:"severity,attr"`
						CweID                                string `xml:"cwe_id,attr"`
						FirstFoundDate                       string `xml:"first_found_date,attr"`
						CveSummary                           string `xml:"cve_summary,attr"`
						SeverityDesc                         string `xml:"severity_desc,attr"`
						Mitigation                           string `xml:"mitigation,attr"`
						VulnerabilityAffectsPolicyCompliance string `xml:"vulnerability_affects_policy_compliance,attr"`
					} `xml:"vulnerability"`
				} `xml:"vulnerabilities"`
				ViolatedPolicyRules string `xml:"violated_policy_rules"`
			} `xml:"component"`
		} `xml:"vulnerable_components"`
	} `xml:"software_composition_analysis"`
}

type Category struct {
	Text         string `xml:",chardata"`
	Categoryid   string `xml:"categoryid,attr"`
	Categoryname string `xml:"categoryname,attr"`
	Pcirelated   string `xml:"pcirelated,attr"`
	Desc         struct {
		Text string `xml:",chardata"`
		Para []struct {
			Text     string `xml:",chardata"`
			AttrText string `xml:"text,attr"`
		} `xml:"para"`
	} `xml:"desc"`
	Recommendations struct {
		Text string `xml:",chardata"`
		Para struct {
			Text     string `xml:",chardata"`
			AttrText string `xml:"text,attr"`
		} `xml:"para"`
	} `xml:"recommendations"`
	Cwe []struct {
		Text        string `xml:",chardata"`
		Cweid       string `xml:"cweid,attr"`
		Cwename     string `xml:"cwename,attr"`
		Pcirelated  string `xml:"pcirelated,attr"`
		Owasp       string `xml:"owasp,attr"`
		Owasp2013   string `xml:"owasp2013,attr"`
		Certjava    string `xml:"certjava,attr"`
		Description struct {
			Chardata string `xml:",chardata"`
			Text     struct {
				Text     string `xml:",chardata"`
				AttrText string `xml:"text,attr"`
			} `xml:"text"`
		} `xml:"description"`
		Staticflaws struct {
			Text string `xml:",chardata"`
			Flaw []struct {
				Text                     string `xml:",chardata"`
				Severity                 string `xml:"severity,attr"`
				Categoryname             string `xml:"categoryname,attr"`
				Count                    string `xml:"count,attr"`
				Issueid                  string `xml:"issueid,attr"`
				Module                   string `xml:"module,attr"`
				Type                     string `xml:"type,attr"`
				Description              string `xml:"description,attr"`
				Note                     string `xml:"note,attr"`
				Cweid                    string `xml:"cweid,attr"`
				Remediationeffort        string `xml:"remediationeffort,attr"`
				ExploitLevel             string `xml:"exploitLevel,attr"`
				Categoryid               string `xml:"categoryid,attr"`
				Pcirelated               string `xml:"pcirelated,attr"`
				DateFirstOccurrence      string `xml:"date_first_occurrence,attr"`
				RemediationStatus        string `xml:"remediation_status,attr"`
				CiaImpact                string `xml:"cia_impact,attr"`
				GracePeriodExpires       string `xml:"grace_period_expires,attr"`
				AffectsPolicyCompliance  string `xml:"affects_policy_compliance,attr"`
				MitigationStatus         string `xml:"mitigation_status,attr"`
				MitigationStatusDesc     string `xml:"mitigation_status_desc,attr"`
				Sourcefile               string `xml:"sourcefile,attr"`
				Line                     string `xml:"line,attr"`
				Sourcefilepath           string `xml:"sourcefilepath,attr"`
				Scope                    string `xml:"scope,attr"`
				Functionprototype        string `xml:"functionprototype,attr"`
				Functionrelativelocation string `xml:"functionrelativelocation,attr"`
				Annotations              struct {
					Text       string `xml:",chardata"`
					Annotation []struct {
						Text        string `xml:",chardata"`
						Action      string `xml:"action,attr"`
						Description string `xml:"description,attr"`
						User        string `xml:"user,attr"`
						Date        string `xml:"date,attr"`
					} `xml:"annotation"`
				} `xml:"annotations"`
			} `xml:"flaw"`
		} `xml:"staticflaws"`
	} `xml:"cwe"`
}

type VeracodeSeverity struct {
	Medium          int
	HighAndVeryHigh int
}

// VeracodeCredentials centralized credentials for veracode account and OpsGenie reporting
type VeracodeCredentials struct {
	Username string
	Password string
}

// Veracode Additional arguments
type VeracodeArgs struct {
	Command   string
	AppID     string
	AppName   string
	BuildID   string
	BuildName string
}

//Find the last build information (build id)
func VeracodeLastBuildInfo(credentials VeracodeCredentials, appid *string, Binfo *BuildInfo) error {
	url := "https://analysiscenter.veracode.com/api/5.0/getbuildinfo.do"
	extraParams := map[string]string{
		"app_id": *appid,
	}

	bodyString, err := makeRequest(url, credentials.Username, credentials.Password, extraParams, REQ_BUILD_INFO)

	if err != nil {
		return err
	}
	err = xml.Unmarshal([]byte(bodyString), &Binfo)
	if err != nil {
		return errors.New(BUILD_NOT_FOUND)
	}
	return nil
}

//Find a specific build information (build id)
func VeracodeBuildInfo(credentials VeracodeCredentials, appid *string, build_id *string, Binfo *BuildInfo) error {
	url := "https://analysiscenter.veracode.com/api/5.0/getbuildinfo.do"
	extraParams := map[string]string{
		"app_id":   *appid,
		"build_id": *build_id,
	}

	bodyString, err := makeRequest(url, credentials.Username, credentials.Password, extraParams, REQ_BUILD_INFO)

	if err != nil {
		return err
	}
	err = xml.Unmarshal([]byte(bodyString), &Binfo)
	if err != nil {
		return errors.New(BUILD_NOT_FOUND)
	}
	return nil
}

func deleteAppLastBuild(credentials VeracodeCredentials, appid string) error {
	extraParams := map[string]string{
		"app_id": appid,
	}

	bodyString, err := makeRequest("https://analysiscenter.veracode.com/api/4.0/deletebuild.do", credentials.Username, credentials.Password, extraParams, REQ_DELETE_BUILD)
	//log.Println(bodyString)
	_ = bodyString
	return err

}

func downloadFullReport(credentials VeracodeCredentials, build_id *string) (Detailedreport, error) {

	var appBuildReport Detailedreport

	url := "https://analysiscenter.veracode.com/api/5.0/detailedreport.do?build_id=" + *build_id
	bodyString, err := makeRequest(url, credentials.Username, credentials.Password, nil, REQ_FULL_REPORT)

	if err != nil {
		return appBuildReport, err
	}
	err = xml.Unmarshal([]byte(bodyString), &appBuildReport)
	if err != nil {
		return appBuildReport, errors.New(bodyString)
	}

	return appBuildReport, err
}

func ScanCheckStatus(build *BuildInfo) error {
	var err error
	if build.AppID != "" {
		log.Printf(BUILD_STATUS, build.Build.AnalysisUnit.Status)
		switch build.Build.AnalysisUnit.Status {
		case STATUS_SCAN_IS_READY:
			err = errors.New(APP_IS_NOT_OK)
		case STATUS_SCAN_IN_PROGRESS:
			err = errors.New(SCAN_IS_IN_PROGRESS)
		case STATUS_SCAN_INCOMPLETE:
			err = errors.New(STATUS_SCAN_INCOMPLETE)
		case STATUS_PRE_SCAN_FAILED:
			err = errors.New(STATUS_PRE_SCAN_FAILED)
		case STATUS_PRE_SCAN_SUCCESS:
			err = errors.New(STATUS_PRE_SCAN_SUCCESS)
		case STATUS_PRE_SCAN_SUBMITTED:
			err = errors.New(STATUS_PRE_SCAN_SUBMITTED)
		case STATUS_NO_MODULES_DEFINED:
			err = errors.New(STATUS_NO_MODULES_DEFINED)
		default:
			err = errors.New(STATUS_SCAN_UNKNOWN + build.Build.AnalysisUnit.Status)
		}
	} else {
		err = errors.New(BUILD_NOT_FOUND)
	}
	return err
}

func SeveritiesNotApproved(report *Detailedreport) (VeracodeSeverity, error) {
	var severity_total VeracodeSeverity
	for _, severity := range report.Severity {
		if severity.Category != nil {
			switch severity.Level {
			case "3":
				SeverityCounter(severity.Category, &severity_total.Medium)
			case "4":
				SeverityCounter(severity.Category, &severity_total.HighAndVeryHigh)
			case "5":
				SeverityCounter(severity.Category, &severity_total.HighAndVeryHigh)
			}
		}
	}
	return severity_total, nil
}

func SeverityCounter(_category []Category, total *int) error {
	for _, category := range _category {
		if category.Cwe != nil {
			for _, cwe := range category.Cwe {
				if cwe.Staticflaws.Flaw != nil {
					for _, flaw := range cwe.Staticflaws.Flaw {
						if flaw.MitigationStatus != "accepted" && flaw.RemediationStatus != "Fixed" {
							*total++
						}
					}
				}
			}
		}
	}
	return nil
}

func FindAppIdByName(credentials VeracodeCredentials, app_name *string, app_id *string) error {
	var appList VeracodeApplist
	url := "https://analysiscenter.veracode.com/api/5.0/getapplist.do"
	bodyString, err := makeRequest(url, credentials.Username, credentials.Password, nil, REQ_FIND_APPID)

	if err != nil {
		return err
	}
	err = xml.Unmarshal([]byte(bodyString), &appList)
	if err != nil {
		log.Println(err)
	}

	for _, item := range appList.App {
		if item.AppName == *app_name {
			*app_id = item.AppID
			log.Println("App ID: ", *app_id)
			return err
		}
	}

	return errors.New(APP_NOT_FOUND)
}

func FindBuildIdByBuildName(credentials VeracodeCredentials, build_id *string, appid *string, build_name *string) error {
	var buildList VeracodeBuildList
	url := "https://analysiscenter.veracode.com/api/5.0/getbuildlist.do"
	extraParams := map[string]string{
		"app_id": *appid,
	}

	bodyString, err := makeRequest(url, credentials.Username, credentials.Password, extraParams, REQ_BUILDLIST)

	if err != nil {
		return err
	}

	err = xml.Unmarshal([]byte(bodyString), &buildList)
	if err != nil {
		log.Println(err)
		return err
	}

	for _, item := range buildList.Build {
		if item.Version == *build_name {
			*build_id = item.BuildID
			log.Println("Build ID: ", *build_id)
			return err
		}
	}
	return errors.New(BUILD_NOT_FOUND)
}
