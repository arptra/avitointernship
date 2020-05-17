from requests import get
from requests import post
from random import randint


def make_req(ip):
    headers = {'X-Forwarded-For': ip}
    r = get('http://localhost:8181', headers=headers)
    return r.status_code


def set_settings(prefix, num_con, limit_time, ban_time, delete_time):
    settings = {'Prefix': str(prefix),
                'NumCon': str(num_con),
                'LimitTime': str(limit_time),
                'BanTime': str(ban_time),
                'DeleteTime': str(delete_time)}
    post('http://localhost:8181/change_settings', data=settings)


def random_ip():
    a = randint(0, 255)
    b = randint(0, 255)
    c = randint(0, 255)
    d = randint(0, 255)
    ip = str(a) + '.' + str(b) + '.' + str(c) + '.' + str(d)
    return ip


def req_for_raise_429(limit, ip_set):
    i = 0
    while i < limit:
        for ip in ip_set:
            status = make_req(ip)
            if status == 429:
                return "OK"
            elif status != 200 and status != 429:
                return "ERROR -> bad answer from server"
            i += 1


def check_after_raise_429(limit, ip_set):
    i = 1
    while i < limit - 1:
        for ip in ip_set:
            if make_req(ip) != 429:
                print("ERROR -> status code from server should be 429")
                return "KO"
            i += 1
            if i == limit:
                break
    return "OK"


def iplist(start_ip, prefix, amount, octet):
    ip_list = []
    if 23 < prefix < 32:
        ip_list = [start_ip + '.' + str(randint(256 - 2 ** (32 - prefix), 255 - 2 ** (32 - prefix - 1))) for _ in range(amount)]
    elif 15 < prefix < 24:
        for _ in range(amount):
            c = randint(int(octet), 255 - 2 ** (24 - prefix - 1))
            d = randint(0, 255)
            ip = start_ip + '.' + str(c) + '.' + str(d)
            ip_list.append(ip)
    elif 7 < prefix < 16:
        for _ in range(amount):
            b = randint(int(octet), 255 - 2 ** (16 - prefix - 1))
            c = randint(0, 255)
            d = randint(0, 255)
            ip = start_ip + '.' + str(b) + '.' + str(c) + '.' + str(d)
            ip_list.append(ip)
    elif 0 < prefix < 8:
        for _ in range(amount):
            a = randint(256 - 2 ** (8 - prefix), 255 - 2 ** (8 - prefix - 1))
            b = randint(0, 255)
            c = randint(0, 255)
            d = randint(0, 255)
            ip = str(a) + '.' + str(b) + '.' + str(c) + '.' + str(d)
            ip_list.append(ip)
    elif prefix == 0:
        for _ in range(amount):
            ip = random_ip()
            ip_list.append(ip)
    return ip_list


def get_iplist(prefix, amount):
    ip = ""
    ip_set = ""
    if 23 < prefix < 32:
        a = randint(0, 255)
        b = randint(0, 255)
        c = randint(0, 255)
        d = 256 - 2 ** (32 - prefix)
        ip = str(a) + '.' + str(b) + '.' + str(c) + '.' + str(d)
        start_ip = str(a) + '.' + str(b) + '.' + str(c)
        ip_set = iplist(start_ip, prefix, amount, d)
    elif 15 < prefix < 24:
        a = randint(0, 255)
        b = randint(0, 255)
        c = randint(256 - 2 ** (24 - prefix), 255)
        ip = str(a) + '.' + str(b) + '.' + str(c) + '.' + str(0)
        start_ip = str(a) + '.' + str(b)
        ip_set = iplist(start_ip, prefix, amount, c)
    elif 7 < prefix < 16:
        a = randint(0, 255)
        b = randint(256 - 2 ** (16 - prefix), 255)
        ip = str(a) + '.' + str(b) + '.' + str(0) + '.' + str(0)
        start_ip = str(a)
        ip_set = iplist(start_ip, prefix, amount, b)
    elif 0 < prefix < 8:
        a = randint(256 - 2 ** (8 - prefix), 255)
        ip = str(a) + '.' + str(0) + '.' + str(0) + '.' + str(0)
        start_ip = str(a)
        ip_set = iplist(start_ip, prefix, amount, a)
    elif prefix == 32:
        ip = "255.255.255.255"
        ip_set = [ip]
    elif prefix == 0:
        ip = "0.0.0.0"
        ip_set = iplist(ip, prefix, amount, 0)
    return ip, ip_set


def generate(start_ip, prefix, amount):
    if 23 < prefix < 32:
        octet = 256 - 2 ** (32 - prefix)
        ipset = iplist(start_ip, prefix, amount, octet)
        return ipset
    elif 15 < prefix < 24:
        octet = 256 - 2 ** (24 - prefix)
        ipset = iplist(start_ip, prefix, amount, octet)
        return ipset
    elif 7 < prefix < 16:
        octet = 256 - 2 ** (16 - prefix)
        ipset = iplist(start_ip, prefix, amount, octet)
        return ipset
    elif 0 < prefix < 8:
        octet = 256 - 2 ** (8 - prefix)
        ipset = iplist(start_ip, prefix, amount, octet)
        return ipset


def func_test(limit, amount):
    prefix_set = range(0, 32)[::-1]
    for prefix in prefix_set:
        if 23 < prefix < 32:
            set_settings(prefix, limit, 1, 1, 1)
            start_ip = str(128) + '.' + str(1) + '.' + str(1)
            ip_set = generate(start_ip, prefix, amount)
            #[print(item) for item in ip_set]
            print("subnet: ", start_ip,
                  "PREFIX = ", prefix,
                  "TEST_RAISE_429_STATUS = ", req_for_raise_429(limit, ip_set),
                  "TEST_CHECK_429_STATUS = ", check_after_raise_429(limit, ip_set))
        elif 15 < prefix < 24:
            set_settings(prefix, limit, 1, 1, 1)
            start_ip = str(128) + '.' + str(1)
            ip_set = generate(start_ip, prefix, amount)
            #[print(item) for item in ip_set]
            print("subnet: ", start_ip,
                  "PREFIX = ", prefix,
                  "TEST_RAISE_429_STATUS = ", req_for_raise_429(limit, ip_set),
                  "TEST_CHECK_429_STATUS = ", check_after_raise_429(limit, ip_set))
        elif 7 < prefix < 16:
            set_settings(prefix, limit, 1, 1, 1)
            start_ip = str(128)
            ip_set = generate(start_ip, prefix, amount)
            #[print(item) for item in ip_set]
            print("subnet: ", start_ip,
                  "PREFIX = ", prefix,
                  "TEST_RAISE_429_STATUS = ", req_for_raise_429(limit, ip_set),
                  "TEST_CHECK_429_STATUS = ", check_after_raise_429(limit, ip_set))
        elif 0 < prefix < 8:
            set_settings(prefix, limit, 1, 1, 1)
            start_ip = "128"
            ip_set = generate(start_ip, prefix, amount)
            #[print(item) for item in ip_set]
            print("subnet: ", start_ip,
                  "PREFIX = ", prefix,
                  "TEST_RAISE_429_STATUS = ", req_for_raise_429(limit, ip_set),
                  "TEST_CHECK_429_STATUS = ", check_after_raise_429(limit, ip_set))


def core_test(prefix, amount, limit):
    ip, ipset = get_iplist(prefix, amount)
    print(req_for_raise_429(limit, ipset))
    print(check_after_raise_429(limit, ipset))


if __name__ == "__main__":
    func_test(50, 10)


