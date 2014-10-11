tc-swas
=======

##Challenge Overview

Welcome to the simple web API server challenge ! 
This is the first challenge of our second stage ‘Develop Backend Services with Golang’. The Goal of this challenge is to develop a simple web API server with Golang. The server provides an API to authenticate users. 
This time we will evaluate your submission based on the scorecard. We totally recommend you to read the articles mentioned in ‘Final Submission Guidelines - Code Guidelines’ section.  The person who gets the highest score wins. In case of a tie, the person to submit earlier wins.
If you have any questions, ask and get clarification in the forum.

##API Spec
This API is to authenticate user for a domain by username and password via HTTP. The domain name is included as a part of the endpoint.

###Endpoint
`/api/2/domains/{domain name}/proxyauth`
 We use port 80 but we would like to use other ports such as 8080 for testing.

###Request
####Request Method
`POST`

####Parameters
* `username`
* `password`
‘password’ parameter is encrypted with the following logic
‘{SHA256}’ + Base64 encoded SHA256 digest of the user’s password
Example
```
original password : abcd1234
password parameter : {SHA256}6c7nGrky_ehjM40Ivk3p3-OeoEm9r7NCzmWexUULaa4=
```

####ContentType
`application/x-www-form-urlencoded`

####Sample
Request parameters
```
domain name : topcoder.com
username : takumi
password : {SHA256}2QJwb00iyNaZbsEbjYHUTTLyvRwkJZTt8yrj4qHWBTU=
Original password is ‘ilovego’
```
Request to a server running on localhost with cURL
`curl --data "username=takumi&password={SHA256}2QJwb00iyNaZbsEbjYHUTTLyvRwkJZTt8yrj4qHWBTU=" http://localhost/api/2/domains/topcoder.com/proxyauth`

###Response
####StatusCode
Use 200 to indicate that the request is processed successfully. Even if we get some application errors such as ‘password unmatch’ or validation errors of parameters, status code should be 200. 404 is used when the domain name is not supported. 500 is used for system errors.
200    Successfully processed the request
404    No such domain
500    Server error

####Format
Return JSON data for status code 200.
In case of success.
```
{ 
        "access_granted": true
}
```
In case of authentication failure or validation errors. The 'reason' is always same.
```
{
        "access_granted": false, "reason": "denied by policy"
}
```
No data should be returned for status code 404 and 500.

####ContentType
`application/json`

###Authentication Logic
This time we use a json file attached (users.json) for data store. 
When you receive a request to appirio.com domain with username ‘jun’ and password, you are supposed to find a record for jun under appirio.com domain in users.json. Encrypt jun’s password you get from the json file, then compare the encrypted password and the password received. If they are same, the authentication succeeds.

###Note
No need to handle signals for this challenge

###Test
Prepare your test script to cover the following cases.
```
Case1 Success
topcoder.com domain
StatusCode 200

Case2 Success
appirio.com domain
StatusCode 200

Case3 Failure
password unmatch
StatusCode 200

Case4 Failure
username not found
StatusCode 200

Case5 Failure
domain not found 
StatusCode 404
```

##Final Submission Guidelines

###Code Guidelines
Follow the practices mentioned in the articles below.
http://golang.org/doc/effective_go.html
https://code.google.com/p/go-wiki/wiki/CodeReviewComments#Go_Code_Review_Comments
We have one note specific to this challenge.
Use  ‘lower_case_with_underscore’ name for package, file or directory. However, try to avoid underscores and prefer short names
Submission Deliverables

##Source code
Format your code with ‘gofmt’ command.
Test script that covers the test cases
Simple README to explain your deliverables

##External Libraries
We believe this challenge is not so complicated that we can complete without any external libraries. But if you would like to use external libraries please follow the guidelines below.
Do not use libraries developed with languages other than Golang
Do not use GPL libraries and LGPL libraries
MIT, Apache and BSD libraries are available
Please mention about external libraries you used in your README
