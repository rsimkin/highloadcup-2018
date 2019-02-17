FROM centos:7

WORKDIR /root
RUN mkdir cup_go

RUN yum install -y wget && \
    wget https://storage.googleapis.com/golang/go1.11.4.linux-amd64.tar.gz

# Устанавливаем Go, создаем workspace и папку проекта
RUN tar -C /usr/local -xzf go1.11.4.linux-amd64.tar.gz && \
    mkdir go && mkdir /root/cup_go/src && mkdir /root/cup_go/bin && mkdir /root/cup_go/pkg

# Задаем переменные окружения для работы Go
ENV PATH=${PATH}:/usr/local/go/bin GOROOT=/usr/local/go GOPATH=/root/cup_go/src GOBIN=/root/cup_go/bin
RUN yum install -y git
RUN go get -u github.com/valyala/fasthttp
RUN go get github.com/qiangxue/fasthttp-routing
RUN ulimit -n 65535
# Копируем наш исходный main.go внутрь контейнера, в папку go/src/dumb
COPY src/ /root/cup_go/src/
# Компилируем и устанавливаем наш сервер
RUN cd /root/cup_go/src/ && go install

# Открываем 80-й порт наружу
EXPOSE 80

# Запускаем наш сервер
CMD ./cup_go/bin/src