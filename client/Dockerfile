FROM python:3.7
RUN pip3 install --upgrade pip
RUN pip3 install requests
ADD . /client/
WORKDIR /client/
#CMD python ext_test/tester.py