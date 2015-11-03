# susi-restapi
This is a simple REST endpoint to event mapper, written in golang
## Installation
To install just use the go command:
```
go get github.com/webvariants/susi-gowebstack
```
This will install the REST server with its dependencies to your $GOROOT
## Usage
The susi-restapi binary need some commandline options
```
Usage of ./susi-restapi:
  -cert="cert.pem": certificate to use
  -https=false: whether to use https or not
  -key="key.pem": key to use
  -mapping="endpoints.json": the endpoint-event mapping file
  -susiaddr="localhost:4000": susiaddr to use
  -webaddr=":8080": webaddr to use
```

susiaddr, key and cert are manadatory to communicate to your local susi-core instance (See the [Susi-Readme](https://github.com/webvariants/susi) to get an susi-core instance running)

You can specify your REST endpoints like this:
```json
{
    "/api/v1/flower/{id}": {
        "GET": "flower::get",
        "POST": "flower::add",
        "DELETE": "flower::remove",
        "PATCH": "flower::update"
    }
}
```
Parameters in the path are passed as payload to the related event. POST, PUT and PATCH also copy their (json) body to the payload.

## Contributing
1. Fork it!
2. Create your feature branch: `git checkout -b my-new-feature`
3. Commit your changes: `git commit -am 'Add some feature'`
4. Push to the branch: `git push origin my-new-feature`
5. Submit a pull request :D

## License
MIT License -> feel free to use it for any purpose!

See [LICENSE](LICENSE) file
