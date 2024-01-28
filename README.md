# fs-chomp

## Installation
## 1. Frontend
- cd ..\frontend\
- npm run build:win
- npm run build:macos
- npm run build:linux
## 2. Backend
  - cd backend
  - go mod tidy
  - go build -o ..\frontend\dist\win-unpacked\
## 3. run
  - cd .\dist\win-unpacked\
  - move pkl files, json files, and py file to this folder 
  - frontend.exe
## 4. to run for dev
  - cd backend
  - go build -o ..\frontend\
  - cd ..\frontend\
  - backend
  - npm run dev