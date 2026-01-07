This Go application takes a directory path and a string (pattern to be matched) as CLI input
and searches recursivel for the pattern in all the files that fall in the directory path.
This search is performed by multiple goroutines.

Here is how it works:

First it creates N goroutines where N is the total no. of cores in the system. Since there
"grep" is a CPU-intensive job, its not beneficial to create goroutines more than the no.
of cores.

The main() performs the directory traversal and places each file(the path of the file)
on a channel (a channel common to all the goroutines). Each goroutine would pick a file
as and when it finishes its previous file. It would read the file line by line and perform
"grep" activity.

Here synchronization must be performed between main() and goroutines. This synchroniztion
has be performed with the help of wait group and and atomic boolean variable.
