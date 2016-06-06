# Central oauth2 sso for services

## Config file:

``` 
{
	"port":"8080",
    "secure":false,
    "tls_key":"",
    "tls_cert":"",
	"cookie_name_space":"instamojo",
	"cookie_secret":"",
	"google_client_id":"google client id",
	"google_client_secret":"google client secret",
    "google_auth_scope":"openid profile email",
    "google_domain":"instamojo.com"
}```

pass the path of the config file as command line argument like this - ./oauth2_central -config-file=path/to/file
if none is passed, program will look for config_file.json in the project root.


## Test and Install
Do a make all to build the project

