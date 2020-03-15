# Simpleton

Simpleton is a dead simple UDP to database logging solution that just
accepts UDP packets and stores them in an SQLite3 database.  Per
default it stores them in `simpleton.db` in the current directory, but
you can override this with command line options.

This isn't terribly useful for anything but really simple testing, but
you can expand on it.

## Building

In order to build simpleton you just run `make` and the binary will
turn up in the `bin` directory.  Per default it will build for OSX.

    make

## Building for other platforms

To build for other platforms please edit the `GOOS` and `GOARCH`
variables in the `Makefile`.  You can also enter these parameters on
the command line when you run `make`, like this.

    GOOS=linux GOARCH=amd64 make
	
You can find the values for these variables for different platforms in 
[syslist.go](https://github.com/golang/go/blob/master/src/go/build/syslist.go),
but the most common values are:

| OS      | GOOS     | GOARCH |
|---------|----------|--------|
| OSX     | darwin   | amd64  |
| Linux   | linux    | amd64  |
| Windows | windows  | amd64  |

Of course, you can cross compile (eg compile Linux binaries on OSX
machines) by just setting the right combination of GOOS and GOARCH,
though Windows you might run into trouble.  (I haven't built this for
Windows).
  
## Running 

The binary will turn up in `bin`, so you can run it from the main
directory with:

    bin/simpleton
	
To list the command line options, you use the `-h` flag:

    bin/simpleton -h
	
Here is an example of running Simpleton with options to make it listen
to a particular interface (10.1.0.3 in the example) and port (7788)
and store the database in `/tmp/simpleton.db`:

    bin/simpleton -u 10.1.0.3:7788 -d /tmp/simpleton.db
	
## Poking around the database

If you want to poke around the resulting database you can install
SQLite3 on your machine and inspect the database using the `sqlite3`
command.  To open the database in the previous example just run:

    sqlite3 /tmp/simpleton.db
	
Type `.schema` to see the very simple database schema.  You can now
perform SQL statements on the data.

Note: I'm not entirely sure about the concurrency of SQLite3 so I
wouldn't use the database as an integration point (you never should).

This is also why the code has a mutex lock around database accesses.
The code was taken from a project that has multiple goroutines
accessing the database.  This program doesn't have that, but I left
the mutex locking in just as a reminder.

For production uses you should use a PostgreSQL database or similar,
that is built for concurrency.  But for small experiments and when you
have limited concurrency, SQLite3 is a surprisingly capable little
beast.


## Accessing via HTTP interface

Note that the HTTP interface has **no authentication or security
mechanisms** so don't use this for anything other than testing.  The
default address of the web interface is:

    http://localhost:8008/

The web interface is quite simple.  You have two URLs that access
data:

    /data
	/data/{id}
	
The first returns a JSON array, the second only returns the payload of
the data entry given by ID.  The `/data` path will be limited to just
the 20 newest entries in the database, but you can page through the
database by setting `offset` and `limit` URL parameters:

    /data?offset=10&limit=10
	

Simpleton supports having a directory with static files so that you
can make some HTML pages with useful links to the content or perhaps
to host JS-frontend applications.

Check the command line help to see the parameters you can fiddle with.
