echo
echo ---- Case1 Success, topcoder.com domain, StatusCode 200 ----
curl -i --data "username=takumi&password={SHA256}2QJwb00iyNaZbsEbjYHUTTLyvRwkJZTt8yrj4qHWBTU=" http://localhost/api/2/domains/topcoder.com/proxyauth
echo
echo ---- Case2 Success, appirio.com domain, StatusCode 200 ----
curl -i  --data "username=jun&password={SHA256}/Hnfw7FSM40NiUQ8cY2OFKV8ZnXWAvF3U7/lMKDwmso=" http://localhost/api/2/domains/appirio.com/proxyauth
echo
echo ---- Case3 Failure, password unmatch, StatusCode 200 ----
curl -i --data "username=jun&password={SHA256}2QJwb00iyNaZbsEbjYHUTTLyvRwkJZTt8yrj4qHWBTU=" http://localhost/api/2/domains/appirio.com/proxyauth
echo
echo ---- Case4 Failure, username not found, StatusCode 200 ----
curl -i --data "username=bryfry&password={SHA256}2QJwb00iyNaZbsEbjYHUTTLyvRwkJZTt8yrj4qHWBTU=" http://localhost/api/2/domains/topcoder.com/proxyauth
echo
echo ---- Case5 Failure, domain not found, StatusCode 404 ----
curl -i --data "username=takumi&password={SHA256}2QJwb00iyNaZbsEbjYHUTTLyvRwkJZTt8yrj4qHWBTU=" http://localhost/api/2/domains/bryfry.com/proxyauth
