#!/bin/bash

go tool compile -S correct.go > correct.S.txt
go tool compile -S incorrect.go > incorrect.S.txt
go tool compile -S incorrect1stmt.go > incorrect1stmt.S.txt

go build correct.go
go build incorrect.go
go build incorrect1stmt.go

go tool objdump -s main.CorrectCast correct > correct.ELF.txt
go tool objdump -s main.WithVariable incorrect > incorrect.ELF.txt
go tool objdump -s main.WithVariable1Stmt incorrect1stmt > incorrect1stmt.ELF.txt

