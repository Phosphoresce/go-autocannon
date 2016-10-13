Autocannon
==========

- Website: https://github.com/phosphoresce/go-autocannon
- IRC: `TO-DO`

Autocannon is a customizable web application testing tool. It supports cURL syntax and focuses on functionality testing and performance testing.

Documentation
-------------

This is it.

Getting Started
---------------

Autocannon is designed to run within a very minimal docker container so all you need to get started is [Docker](https://www.docker.com).  

Once you have docker installed simply pull the image.  

```sh
$ docker pull phosphoresce/autocannon
```

Running autocannon is controlled through the docker daemon.

```sh
$ docker run phosphoresce/autocannon -h
```

If you would like to run Autocannon without Docker, feel free to clone the repository and build the binary with the included Makefile.  

Thats pretty much it so far!

Developing Autocannon
---------------------

Build the project with the included Makefile and submit me a pull request with any changes.  

WIP Notes:
why dont i just output a whole shit ton of responses to a file?
that way i can then include that file in the next round for a particular field?
