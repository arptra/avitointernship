# Avito internship  
Create an HTTP service that can limit the number of requests from one IPv4 subnet. If there are no restrictions, then you need to produce the same static content.  
[Original test task from Antibot Avito Information Security unit](https://github.com/avito-tech/antibot-developer-trainee)  

## A brief overview of what I did  
#### I implemented:
1) Сlient part on `Python`.  
2) Server part on `Golang`.  
3) Ability to run using `docker-compose up`.  
4) Implemented server reconfiguration `handler` (change prefix, limite, ban, delete time and numder of request).  
5) `handler` to reset the restriction on the subnet address.  
6) `Testing` the functioning of the service.
7) A simple Info webpage to display the necessary information about the service (settings, limited networks, connections), as    well as the ability to reconfigure the service while it run.  
  
![Alt Text](https://github.com/arptra/avitointernship/blob/master/pic/infopage.jpeg)  
  
#### What did I use for this  
1) Ubuntu 18.04.4 LTS (all of my code run and tests on this system).  
2) Golang 1.14 with module github.com/gorilla/mux v1.7.4.  
3) Python 3.6 with module request and tkinter
4) Docker version 19.03.8, build afacb8b7f0
5) Docker-compose version 1.25.3, build d4d1b42b  
  
# How to use  
The fastest way to start using  `docker-compose up`
```
git clone https://github.com/arptra/avitointernship.git  
```
after cloning, move to this directory and run
```
docker-compose up
```
When the containers start, you will see the following on the command line  
  
![Alt Text](https://github.com/arptra/avitointernship/blob/master/pic/docker_compose_up.jpeg)
When the container with the server starts, the container with the client starts and run test  
These are two types of tests in the picture:  
```
1) TEST_RAISE_429  
2) TEST_CHECK_429  
```  
What they mean and how I test can be read more on my [wiki](https://github.com/arptra/avitointernship/wiki)  
If you open your browser(I recommended full screen for right rendering page) and go to http://localhost:8181,  
you will receive a message about the restriction (this is due to testing), you will need to wait a minute,  
then the main page will be available (Info pag)  

![Alt_text](https://github.com/arptra/avitointernship/blob/master/pic/to_many_request.png)  

#### Another way to build an application from source, for this  
You must have installed golange 1.14 and python3.  
Run next command:  
```
git clone https://github.com/arptra/avitointernship.git
cd avitointernship/server
go mod download
go build -o app main.go
./app
```  
With this launch, you can use the flags  
```
./app -p=24 -nc=100 -lt=1 -bt=2 -dt=1
```
```
-p    prefix  in decimal  
-nc   Number of request  
-lt   limite time (time during which it is possible to make a limited number of requests)  
-bt   ban time (how long you will in restricted list)  
-dt   delete time (time after which the routine starts and deletes all obsolete entries in the restricted list)  
  all time is set in minutes  
```  
# Client
After service have started you can run client with visual mode.  
For this you need installed python3 and module request.  
```  
pip3 install request  
```  
```
cd avitointernship/client
```
and run
```
python3 -m ext_test.client
```
![Alt_test](https://github.com/arptra/avitointernship/blob/master/pic/client.png)  

It has two main functions.
1) Generate a list of IP addresses from the same subnet and make a request to the server from each server (first line)
2) Generate a random IP address and make a server request from it (Second line)


