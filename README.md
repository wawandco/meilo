# Meilo
Meilo is an implementation of a test/development SMTP server written in Go. Meilo opens sent emails instead of sending them, 
this allows development workflows where ONLY the SMTP credentials are changed.

## Prerequisites
- Meilo requies Go 1.23 or later, you can download it from [here](https://golang.org/dl/)

- Meilo package:
```bash
go get github.com/wawandco/meilo
```

## Usage
1.  Create a new instance of Meilo

You have the ability to set the options for the SMTP server, these are totally optional, if you don't set them, Meilo will use the default values.

```go
// Start the SMTP server
creds, err := meilo.Start(
	// Directory to put the files
	meilo.WithDir("/my/emails/folder"),

	// Port to use
	meilo.WithPort("1025"),
)
 
if err != nil {
	// Handle the error
	...
}
```
This will start the SMTP server and return the credentials to be used in the email sending process.

2.  Send an email using the SMTP server

```go
//Then you can use the credentials to build the SMTP auth 
auth := smtp.PlainAuth("", creds.User, creds.Password, creds.Host),
body:= []byte("Hello from meilo!")
from := "username@example.com"
to:= []string{"example@example.com"}

// And then the creds instance has an Addr method to use when sending
err = smtp.SendMail(creds.Addr() auth, from ,to, body)
if err != nil {
        // Handle the sending error
        ...

}
```
### Options explained

##### `meilo.WithDir(directory)`
Allows to specify the directory where the emails will be stored, by default it will use the system's temporary directory.

##### `meilo.WithPort(port)`: 
Allows to specify the port of the SMTP server. This is useful when running multiple services in your development environment.


## Roadmap / Ideas
- Web interface
- UI Improvements.
- Deployable service for Staging/Testing.
- Listing historically sent emails.

