# cjpais.com
My website written in Go. This is meant to be nothing more than a fast blogging/diary platform.

This site will contain my random thoughts without the expectation that anyone is going to read them.
This is merely a way for me to log my thoughts.

This project has a few goals in mind:
* To have an infinitely scrolling blog
* This blog has a timeline on the side that can be used to scroll to moments in time (JS). Google Photos like
* There is no database necessary (at the moment)
* Files are always served statically with metadata ingrained into them

A few of these goals are with creating a distributed system in mind that can easily serve static files.
This platform will likely grow over time and will eventually take advantage of a distributed key-value store
to get metadata rather than it being part of the file. I am not sure what this metadata will be at the moment.

Also note there is a missing file from this directory. I may have hard coded my authentication into the site. 
I know this is a absolutely terrible practice, but again I don't want to use a database so a static file 
will suffice as it only needs 1 kv pair. This has been done with scrypt.
