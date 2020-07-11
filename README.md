# Gro data loader for getting data from online banking company Q2.

## Usage

### From go command line
```
go run main.go -s 2020-07-06 -e 2020-07-10
```

### Executable

```
gro_loader.exe -s 2020-07-06 -e 2020-07-10
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
- URL
- User ID
- Process Session ID
- Security Token
- Local File Path

## Command line flags
 - (-c) Config File path (optional). Default value ./config.json
 - (-s) Start Date
 - (-e) End Date
