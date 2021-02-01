# Veracode Command line interface

## Automated way to check application status and DevSecops compliance


Analyze the latest build report according to the Adidas DevSecops compliance.
Checking the latest build status and delete the incomplete scans automatically to make sure the application is ready for a new scan request.

## Requirements
A Veracode API credential with Result API permission.

## Usage

```
Usage of veracode-cli:
  -appname string
        [Mandatory] Application Name
  -buildname string
        [Optional] Build Name
  -command string
		[Mandatory] Veracode command:
				getappid - Finding application ID by using application name.
				getbuildinfo - Showing the information of the lastest build
				getbuildversion - Retrieving the latest build version name
				getbuildid - Retrieving the build id by using build name/version
				buildstatus - Checking the status of the latest build and delete it if it's been stuck
				devsecopscheck - Checking DevSecops check of the latest/specific build.
  -pass string
        [Mandatory] Password
  -user string
        [Mandatory] Username
  -version
        Prints the current version.

```

## Example

#### DevSecops checking of a specific application

```
./veracode-cli -user '<API_KEY>' -pass '<API_SECRET>' -command 'devsecopscheck' -appname '<Application Profile Name>' -buildname <build name>
```

If the application has some none-mitigated flaws with high severity

```
2020/03/10 03:20:29 Veracode CLI version 1.1.0
2020/03/10 03:20:29 DevSecopscheck
2020/03/10 03:20:29 Finding App ID
2020/03/10 03:20:31 App ID:  *****
2020/03/10 03:20:31 Requesting build list
2020/03/10 03:20:31 Build ID:  *****
2020/03/10 03:20:31 Requesting build info
2020/03/10 03:20:32 The build status is [Results Ready]
2020/03/10 03:20:32 Requesting full report
2020/03/10 03:20:35 [Result] { High & V.High: 0 , Medium: 3 }
2020/03/10 03:20:35 [Error] The build has some none-mitigated flaws with high severity!
```

#### Checking Application status.

```
./veracode-cli -user '<API_KEY>' -pass '<API_SECRET>' -command 'buildstatus' -appname '<Application Profile Name>'
```

The latest build will be deleted if the status is incomplete

```
2020/02/03 06:33:01 Veracode CLI version 1.1.0
2020/02/03 06:33:01 Checking Build status
2020/02/03 06:33:01 Finding App ID
2020/02/03 06:33:04 App ID:  ***
2020/02/03 06:33:04 Requesting build info
2020/02/03 06:33:05 The build status is [Incomplete]
2020/02/03 06:33:05 Deleting last build
2020/02/03 06:33:06 [Success] The scan was stuck and successfully deleted.
```

## Cross compile

```
$ go get github.com/mitchellh/gox
$ gox --output veracode-cli_{{.OS}}_{{.Arch}}
```

## Publishing a release

Edit [version](./version) file to change the version accoring to [semver](https://semver.org/). Commit the changes and then create a tag:

```
git tag -a $(make version) -m '$(make version)'
```

Once the tag is created, push the changes and Travis will automatically perform the release.


## License and Software Information

© adidas AG

adidas AG publishes this software and accompanied documentation (if any) subject to the terms of the MIT license with the aim of helping the community with our tools and libraries which we think can be also useful for other people. You will find a copy of the MIT license in the root folder of this package. All rights not explicitly granted to you under the MIT license remain the sole and exclusive property of adidas AG.

NOTICE: The software has been designed solely for the purpose of analyzing the code quality by checking the coding guidelines. The software is NOT designed, tested or verified for productive use whatsoever, nor or for any use related to high risk environments, such as health care, highly or fully autonomous driving, power plants, or other critical infrastructures or services.

If you want to contact adidas regarding the software, you can mail us at _software.engineering@adidas.com_.

For further information open the [adidas terms and conditions](https://github.com/adidas/adidas-contribution-guidelines/wiki/Terms-and-conditions) page.

### License

[MIT](LICENSE)
