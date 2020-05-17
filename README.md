# Avito internship  
Create an HTTP service that can limit the number of requests from one IPv4 subnet. If there are no restrictions, then you need to produce the same static content.  
[Original test task from Antibot Avito Information Security unit](https://github.com/avito-tech/antibot-developer-trainee)  

## A brief overview of what I did  
I implemented:
1) Сlient part on `Python`.  
2) Server part on `Golang`.  
3) Ability to run using `docker-compose up`.  
4) Implemented server reconfiguration `handler` (change prefix, limite, ban, delete time and numder of request).  
5) `handler` to reset the restriction on the subnet address.  
6) `Testing` the functioning of the service.
7) A simple Info webpage to display the necessary information about the service (settings, limited networks, connections), as    well as the ability to reconfigure the service while it run.  
![Alt Text](https://github.com/arptra/avitointernship/blob/master/pic/infopage.jpeg)  
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
1) TEST_RAISE_429  
2) TEST_CHECK_429  
What they mean and how I test can be read more on my [wiki](https://github.com/arptra/avitointernship/wiki)
