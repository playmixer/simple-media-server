go build -o .\build\main.exe ./cmd/app/
rmdir .\build\static\ /s/q
xcopy ".\static\" ".\build\static\" /s/h/e/k/f/c/y
rmdir .\build\templates\ /s/q
xcopy ".\templates\" ".\build\templates\" /s/h/e/k/f/c/y