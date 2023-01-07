# mail-indexer
##### Español [README_es.MD](https://github.com/mata649/mail_indexer/blob/master/README_es.md)
This application is divided into the **Indexer** and **Efinder**. The **indexer** gets emails from a path, then parse the emails to provide them a structure, finally ingesting a zinc engine with them.
**Efinder** is just a web application to access indexed emails, the interface was built in **Vue.js** with **TailwindCs**, and the API was developed in Go using **Chi** as a router, also I took advantage of the **API** to server the build generated by Vue, using the main URL as a **Static File Server**

## Indexer
The Indexer is the core of the application, it receives a path of email as an argument, to ingest the zinc engine with the emails.
### Steps 
1. Receives a path where the emails are
2. Gets a slice with the paths of all the files in the directory
3. Divides the paths into smaller slices, the number of emails per slice can be set in the **config.json** file, with the property emailsPerFile
4. Iterates the divided paths concurrently, the number of goroutines running at the same time can be set in the **config.json**, with the property nWorkers
5. Iterates the paths to get the email from each file, then the email is added to a slice of emails.
6. When all the emails have been obtained, the slice of emails is converted into a buffer of bytes in an NDJSON format
7. Makes the request to the Zinc Engine sending the buffer of bytes as a Binary, the user, password, and host of the Zinc Engine can be set in the **config.json**
### Running the code
#### Important
In Linux, the program throws an error due to the number of files that a process can open, you can avoid this error by changing in the **config.json** the emailsPerFile, by default it is in **1000**, but if you are having problems with this you can set **100** emailsPerFile.

The second option to solve this (**no recommended**) is to increase the number of files allowed to open by a process. You can do it with the command: `ulimit -n 8000`
#### VSCode
For some reason, the previous error just appears if you are running the program directly from a terminal, if you run the code from the integrated terminal in **VSCode** the program does not throw an error even with **1000** emailsPerFiel.

To run the code you can do it with the next command

    ./indexer -emailpath /example/path/to/enron_mail_20110402
Or you can run the main.go file with the next command

    go run cmd/main.go -emailpath /path/example/to/enron_mail_20110402 
There are some flags that you can add if you want to see the profiling of the code

    go run cmd/main.go -emailpath /path/example/to/enron_mail_20110402 -cpuprofile profile/cpu_profile.out -memprofile profile/mem_profile.out 

## Test
The program has some unit tests to test each function, this is helpful when we are doing changes in the code and we want to check if the function has the expected behavior, to run all the tests of the program you can do it with the next command.

    go test ./...
  
## Optimization
The **v1** of the program has some optimization problems, this could be checked in the [profile/v1](https://github.com/mata649/mail_indexer/tree/master/indexer/profile/v2) folder. 
In the **v1** of the program, all are relatively similar until steps 6 and 7 in the **v2** which are :

 6. When all the emails have been obtained, the slice of emails is converted into a buffer of bytes in an NDJSON format
 7. Make the request to the Zinc Engine sending the buffer of bytes as a Binary, the user, password, and host of the Zinc Engine can be set in the **config.json**

In the **v1** after of have been obtained the slice of emails, was saved in a file in a data directory, with a folder with the date-time of the running as name: `2023-01-07 14:13:56`. In that folder, the files were enumerated with a number, for example: **file1.ndjson**, **file2.ndjson**, ...., with the number of emails setted in **emailsPerFile** in each file. Then when we are going to ingest the **Zinc Engine** with the emails, we have to get the paths of the files and read them to get the bytes of buffer to send to the Zinc Engine.

This has some advantages and disadvantages,
#### Advantages

 - We have a historical record of the emails used in each running
 - If some error occurs in the ingestion we can log the file we got the error to check it
 - If some error occurs in the ingestion and we can't continue ingesting more files, we can see what was the last file to be ingested and then continue the ingestion from that file.
#### Disadvantages
 - For each file, we have to open two buffers, the first when we are creating the file and after when we have to read it to ingest the Zinc Engine. 

The difference between creating the file to read it after and sending the buffer directly was from **2m 55s** to **2m 21s**. To be honest, depends on the needs of a real project, if you need an **Indexer** "fault-tolerant" when if some problem occurs you have a checkpoint to continue after the ingestion you can use the **v1**, but if you want an **Indexer** that ingest the **Zinc Engine** fastly, you can use the **v2**. 

Maybe in the **v3** of the program could send the buffer of bytes directly to make the request and also save the file, so we are using the same buffer for both purposes. 
## Efinder
Efinder is a web application to search for a term in indexed emails. It was built with **Vue.js** and **Tailwind.css** to develop the interface, and Go with the **Chi** router to request the **Zinc Engine** and also to server the distribution of the Vue Application, working as a **Static File Server*

You can run the application with the next command:

    ./efinder -host http://localhost:4080 -user admin -password Complexpass#123
### Flags

 - **host**: host of the **Zinc Engine** to make the requests
 - **user**: user needed to do the basic authentication in the **Zinc Engine**
 - **password**: password needed to do the basic authentication in the **Zinc Engine**
