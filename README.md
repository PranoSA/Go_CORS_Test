# CORS Tester


Tests Whether or Not a Request Will Generate a Preflight Request or is Subject To Simple CORS
If a Simple CORS is requested, will test if the request qualified for Simple CORS, Then Tests Whether Simple CORS works

If Simple CORS is not requested, Will Test if Requests passes preflight requests 
Based on Requested Methods, Headers, and Whether Credentials are Used (Cookie, Authorization header)

Remember that the CORS Response is configured per Originating Domain and Per Method, *BUT* , The server or functionality or method
to handle CORS requests for each endpoint on a server may differ, so remember to specify the PATH in the destination path.


## Command Line Options 

### Must :

Specify Origin with (--origin )

Specify Dest with (--dest) [ex: "GET", "GET,POST", "GET,POST,DELETE"]

Specify Methods with (--methods)

## Optional:

Specify Simple CORS (--simple)

Specify Headers (--headers)

Specify Credentials (--credentials)