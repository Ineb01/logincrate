FROM golang:1.19 AS go-build

WORKDIR /logincrate

COPY public.env ./
RUN export $(xargs < public.env)

COPY golang/go.mod golang/go.sum ./
RUN go mod download

COPY golang/*.go ./

RUN go build -o /bin/logincrate

#FROM scratch
#COPY --from=go-build /bin/logincrate /bin/logincrate
CMD ["/bin/logincrate"]