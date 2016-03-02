Notes:
1. All of configuration located in config folder
2. A new file will be written after X (minutes) duration set in conf.json

The parameter is as follows

{
  "Target": "http://mybackend" -> the url and port of your backend server. change this accordingly
  "ListenOn": "127.0.0.1:8888", -> This is the url and port we exposed to the world change this accordingly
  "StatFolder" : "/log/", -> location of log file, make sure it exist and we can write in it. change this accordingly
  "CreateNewFileEveryXMinutes" : 1 -> change this accordingly
}
