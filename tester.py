import requests
import ipaddress
from random import randint


def makeReq(ip):
    headers = {'X-Forwarded-For': ip}
    r = requests.get('http://localhost:8181', headers=headers)
    return (r.status_code)


def randomIp():
    a = randint(0, 255)
    b = randint(0, 255)
    c = randint(0, 255)
    d = randint(0, 255)
    ip = str(a) + '.' + str(b) + '.' + str(c) + '.' + str(d)
    return ip


def randomIpSet(limit, prefix, startIp):
    subnet1 = ipaddress.ip_network(startIp + '/' + str(prefix))
    addrspace = list(subnet1.hosts())
    numOfAddr = len(list(subnet1.hosts()))
    ip_list = [str(addrspace[randint(0, numOfAddr - 1)]) for _ in range(limit)]
    return ip_list


def reqForRaise429(limit, ipSet):
    i = 1
    while i < limit - 1:
        for ip in ipSet:
            if (makeReq(ip) != 200):
                print("ERROR -> status code from server should be 200")
                print(i)
                return "KO"
            i += 1
            if i == limit:
                break
    return "OK"


def checkAfterRaise429(limit, ipSet):
    i = 1
    while i < limit - 1:
        for ip in ipSet:
            if (makeReq(ip) != 429):
                print("ERROR -> status code from server should be 429")
                return "KO"
            i += 1
            if i == limit:
                break
    return "OK"


if __name__ == "__main__":
    ipSet = randomIpSet(25, 24, '123.45.67.0')
    print(reqForRaise429(100, ipSet))
    print(checkAfterRaise429(100, ipSet))
    print(makeReq(randomIp()))
    [makeReq(randomIp()) for _ in range(100)]
