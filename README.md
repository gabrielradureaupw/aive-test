# Aive Test

## The Subject

> "Covid vaccination"

Code the API allowing the following use-cases:
- Make an appointment for the vaccination
  - List the available slots in the different vaccination centers
  - Make the appointment and validate via email
- List the daily appointments for a vaccination centers via an authenticated endpoint

## Guidelines
- Use Golang
- Data must be persisted
- Push your code to a Github repository
- Document what you’ve done

## Evaluation
- Quality of code 
- Right choice about the components
- Usage of the language good practices

## Implementation

```sh
go run .

# if you have the jq command append '|jq' to each curl

curl localhost:8080/appointments # list available time slots

curl localhost:8080/appointments/daily \
-H "Authorization: Basic $(echo -n aive:test | base64)" # list daily appointments

curl localhost:8080/appointments -d '{
  "vaccinationCenter":{
   "city":"",
   "name":""
  },
  "timeSlot":"",
  "email":""
}' # make appointment

# if you click on the link displayed in the logs,
# delete the trailing %5C (url encoded '\') suffix
# (as in 'http://localhost:8080/appointments/confirm?id=1&email=test@aive.io\')

curl localhost:8080/appointments/confirm?id=&email= # confirm email
```

## Step by Step

- start server
  ```sh
  cd aivetest
  go run . # start server
  ```
- no appointment made  
  `
  curl localhost:8080/appointments/daily -H "Authorization: Basic $(echo -n aive:test | base64)" | jq
  `
  ```json
  {
    "dailyAppointments": {}
  }
  ```
- list available slots  
  `curl localhost:8080/appointments | jq | head`
  ```json
  {
  "vaccinationCenters": [
    {
      "city": "Northleach",
      "name": "Jawthorn",
      "slots": [
        "2021-06-23T10:00",
        "2021-06-23T11:00",
        "2021-06-23T12:00",
        "2021-06-23T14:00",
  //... truncated
  ```
- make appointment
  ```sh
  curl localhost:8080/appointments -d '{
    "vaccinationCenter":{
    "city":"Northleach",
    "name":"Jawthorn"
    },
    "timeSlot":"2021-06-23T11:00",
    "email":"test@aive.io"
  }' -i
  HTTP/1.1 202 Accepted
  Date: Tue, 22 Jun 2021 17:22:41 GMT
  Content-Length: 0

  # EMAIL LOGS
  {"level":"info","time":"2021-06-22T19:22:41+02:00","message":"Sending Confirmez votre rendez-vous to test@aive.io with content: '\n\tBlablabla\n\tBla blabla\n\t<a href=\"http://localhost:8080/appointments/confirm?id=1&email=test@aive.io\">Confirmez</a>\n\tFormule de politesse\n\t'"}
  ```
- one appointment made but not confirmed
  `
  curl localhost:8080/appointments/daily -H "Authorization: Basic $(echo -n aive:test | base64)" | jq
  `
  ```json
  {
    "dailyAppointments": {
    "2021-06-23T11:00": [
      {
        "vaccinationCenter": {
          "city": "Northleach",
          "name": "Jawthorn"
        },
        "timeSlot": "2021-06-23T11:00",
        "email": "test@aive.io",
        "confirmed": false
      }
    ]
  }
  ```
- visit `http://localhost:8080/appointments/confirm?id=1&email=test@aive.io` or use curl
  ```json
  {
    "message": "RDV confirmé"
  }
  ```
- one appointment made and confirmed
  `
  curl localhost:8080/appointments/daily -H "Authorization: Basic $(echo -n aive:test | base64)" | jq
  `
  ```json
  {
    "dailyAppointments": {
    "2021-06-23T11:00": [
      {
        "vaccinationCenter": {
          "city": "Northleach",
          "name": "Jawthorn"
        },
        "timeSlot": "2021-06-23T11:00",
        "email": "test@aive.io",
        "confirmed": true
      }
    ]
  }
  ```

### Amount of time taken: around 6 hours