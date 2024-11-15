# ssh-ed25519 key gen

This program allows you to specify a list of keywords
to search for, and it will generate a ton of `ssh-ed25519` keys
and log any matches.

## Installation

Make sure to have git and golang version 1.23.2 or higher installed on your computer.

```shell
$ git clone https://github.com/kyren223/ed25519-key-gen
$ cd ed25519-key-gen
$ go build main.go
```

## Usage

Run the program with a list of keywords

```
$ ./ed25519-key-gen keyword1 keyword2 keyword3
```

Alternatively you can create a `input.txt` file in the directory
where the program runs, and put 1 keyword per line.

To stop the program press Ctrl+C or trigger a SIGINT interrupt.

The keyword and public key will be displayed in STDOUT
The keyword, public key and private key will be appeneded to the `output.txt` file.

The `output.txt` format contains each match in a separate line,
each line has 4 space-separated items: keyword, ssh public key, base64 encoded private key, base64 encoded public key.

## Changing amount of goroutines (concurrent searches)

By default there are 50 goroutines, which are light weight threads (also known as green threads),
if you wish to change it, you can edit `const goroutines = 50` in `main.go` to your desired amount.
