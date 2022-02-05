# Project Title

## FFProbe as a Service

### The What

Do you need to get the metadata/information of your file? You can use this service to do that.
### The Why

There are little ways of getting the metadata of media files without resorting to shipping FFMPEG with your app. This is a standalone service that does that. Drop a URL of the file and get all the details without downloading anything.


## Screenshots/Demo


## Built With

>List of tech languages, frameworks/libraries, and tools used

- [Golang](https://go.dev/)
- [Iris](https://github.com/kataras/iris)
- [Redis](https://redis.io/)
- [RabbitMQ](https://www.rabbitmq.com/)

## Features

[x] Small footprint.
[x] Fast.
[x] can switch between local and a remote worker instance.
[x] Open-source


## Quicker Start
Check out [Metaworka](https://metaworka.herokuapp.com/) and [Durator](https://durator.web.app/) for a free and hosted version of this project.

## Quick Start
While in the root directory, copy and rename the `.env.example` file to `.env` and populate it appropriately. To get the server started, run
```
go get
make serve
```
If you set `QUEUE_REQUEST` to `false` in the `.env` file, you would need to start a worker instance by running.

```
make run
```



## API Reference/Documentation


## Contributing

Issues and pull requests are welcome at [cdEnv](https://github.com/teezzan/cdEnv). This project is intended to be safe, welcoming, and open for collaboration. Users are expected to adhere to the [Contributor Covenant code of conduct](https://www.contributor-covenant.org/version/2/0/code_of_conduct/). We are all human.

## Authors

**[Taiwo Yusuf](https://github.com/teezzan/)**


## Acknowledgments

**[TemiTayo Ogunsusi](https://www.linkedin.com/in/temitayo-ogunsusi)** for building [Metaworka](https://metaworka.herokuapp.com/)
**[Abdullah AbdulFatah](https://www.linkedin.com/in/abdullah-abdulfatah-125189209/)** for building [Durator](https://durator.web.app/)
**[Meg Gutshall](https://github.com/meg-gutshall/)** for her README template. Helped a lot.

## License
This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.