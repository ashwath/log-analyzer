# Log Analyzer
REST Service for log analysis 
- keyword search
- tail logs

## Dependencies
* Docker
  * install using: https://docs.docker.com/get-docker/

## Usage
Checkout the Makefile
```bash
$ make {command}
``` 

## API definitions

### Response Codes
```
200: Success
400: Bad request
50X: Server Error
```

### Example Error Messages
```json
{
  "status": 400,
  "name": "Bad Request",
  "message": "invalid input for {limit}",
  "internal_message": "invalid input for {limit}"
}
```

```json
{
  "status": 500,
  "name": "Internal Server Error",
  "message": "open /var/log/sam.log: no such file or directory",
  "internal_message": "open /var/log/sam.log: no such file or directory"
}
```

## Tail API
### URL 
`v1/logs/tail`

### Method 
`GET`

### URL Params

#### Required
* Log file name: `file_name=[alphanumeric]`

#### Optional
* Last n log enteries, defaults to 20: `lastN=[integer]`

#### Notes
* Results are returned in reverse time ordered

### Examples
**Request:**
```
URL: 'http://localhost:4200/v1/logs/tail?file_name=sample.log&lastN=5' 
```
**Successful Response:**
```json
{
  "results": [
    {
      "file_path": "sample.log",
      "log_entries": [
        "",
        "2012-02-03 20:11:56 SampleClass3 [INFO] everything normal for id 530537821",
        "2012-02-03 20:11:56 SampleClass3 [TRACE] verbose detail for id 1718828806",
        "2012-02-03 20:11:56 SampleClass8 [DEBUG] detail for id 2083681507",
        "2012-02-03 20:11:56 SampleClass7 [TRACE] verbose detail for id 1560323914"
      ]
    }
  ],
  "response_metadata": {
    "next_file": "",
    "next_cursor": 0
  }
} 
```

**Request:**
```
URL: 'http://localhost:4200/v1/logs/tail?file_name=xxxx' 
```
**Error Response:**
```json
{
  "status": 500,
  "name": "Internal Server Error",
  "message": "open /var/log/xxxx: no such file or directory",
  "internal_message": "open /var/log/xxxx: no such file or directory"
}
```


## Search API
### URL
`v1/logs/search`

### Method
`GET`

### URL Params
#### Optional
* Log file name, defaults to all log files: `file_name=[alphanumeric]`
* Search Keyword, defaults to all log enteries: `keyword=[alphanumeric]`
* Number of log entries, defaults to 20: `limit=[integer]`
* Paging metadata: `next_cursor=[integer]` & `next_file=[aphanumeric]`

#### Notes
* Results are returned in reverse time ordered

### Examples
**Request #1:**
```
URL: 'http://localhost:4200/v1/logs/search?limit=5&file_name=sample.log&keyword=5305' 
```
**Successful Response #1:**
```json
{
  "results": [
    {
      "file_path": "/var/log/sample.log",
      "log_entries": [
        "2012-02-03 20:11:56 SampleClass3 [INFO] everything normal for id 530537821",
        "2012-02-03 20:11:35 SampleClass4 [TRACE] verbose detail for id 1530516857",
        "2012-02-03 20:04:31 SampleClass9 [TRACE] verbose detail for id 763035305"
      ]
    }
  ],
  "response_metadata": {
    "next_file": "",
    "next_cursor": 0
  }
} 
```

**Request#2:**
```
URL: 'http://localhost:4200/v1/logs/search?limit=5&file_name=sample.log&next_cursor=98978&next_file=/var/log/sample.log' 
```
**Successful Response #2:**
```json
{
  "results": [
    {
      "file_path": "/var/log/sample.log",
      "log_entries": [
        "",
        "2012-02-03 20:11:56 SampleClass5 [TRACE] verbose detail for id 990982084",
        "2012-02-03 20:11:56 SampleClass6 [DEBUG] detail for id 1546542023",
        "2012-02-03 20:11:56 SampleClass7 [TRACE] verbose detail for id 2067347557",
        "2012-02-03 20:11:56 SampleClass2 [TRACE] verbose detail for id 1414320436"
      ]
    }
  ],
  "response_metadata": {
    "next_file": "/var/log/sample.log",
    "next_cursor": 98687
  }
}
```
