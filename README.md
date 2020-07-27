# go-utils

## What the heck is this?

This is a collection of stuff I use across various of my go projects.

### errors.go

This contains stuff to make error handling easier. This also contains
stuff to build up an error stack that shows the origin of the error.

### logger.go

This is a work in progress to have a standardized way to log stuff,
is able to either log to the console but is mainly designed to send
the logs to a remote address running a logging service.

The library handles the logs using a dedicated go routine. In the
future this should be a worker pool with either a static, or
a dynamically controlled number of workers.