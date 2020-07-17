# Gro data loader for getting data from online banking company Q2.

## Usage

### Parameters
Only -s start (date or date time) and -e end (date or date time) parameters are required.

To view all available flags use:

```
gro_loader.exe help
```

### Run from go command line
```
go run main.go -s 2020-07-06 -e 2020-07-10
```

### Executable

#### Dates only
```
gro_loader.exe -s 2020-07-06 -e 2020-07-10

```

#### Date Time
```
gro_loader.exe -s 2020-07-06T08:00:00 -e 2020-07-10T23:59:59
```

## Request Required Parameters

### Header
* Content-Type - application/json
* Content-Length
* Host

### Body
```json
{
  "messageContext": {
    "fiId": "999999",
    "userList": {
      "user": [
        {
          "userId": "???",
          "processorSessionId": "???",
          "userIdType": "FI_USER_ID"
        }
      ]
    },
    "customData": {
      "valuePair": [
        {
          "name": "securityToken",
          "value": "???"
        }
      ]
    }
  },
  "applicationFilter": {
    "applicationIdList": {
      "applicationId": []
    },
    "applicationDateRange": {
      "startDateTime": "2020-07-11T00:00:00",
      "endDateTime": "2020-07-11T23:59:59"
    }
  }
}
```

## Config information
- User ID
- Process Session ID
- Security Token
- Local File Path

## Command line flags

```
  -c string
        Configuration file path. (default "./config.json")

  -e string
        End DateTime in 'YYYY-MM-DD' or 'YYYY-MM-DDTHH-MM-SS' format. (Required)
  -o string
        Output file path (default "./data.json")
  -pr int
        The number of pages requested; this is the number that should be changed in
        the case there are more than 10 queries in a search your financial institution
        is trying to utilize. For example, if the financial institution has 300
        applications in the month of December (using the same time constraints from
        before), then the URL parameter “pr” would be set to any value from 1-30
        (taking the applications 10 entries at a time – 1-10, 11-20, 21-30, etc). (default 1)
  -ps int
        Stands for the size of the page.
        Since the standard is having 10 applications
        maximum per page (as seen in the admin portal search queries), we recommend
        only using a maximum of 10 for this value in the URL. (default 10)
  -s string
        Start DateTime in 'YYYY-MM-DD' or 'YYYY-MM-DDTHH-MM-SS' format. (Required)
```
