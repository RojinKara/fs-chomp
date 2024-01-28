# fs-chomp

## Installation
## 1. Frontend
- cd frontend
- npm run build:win
- npm run build:macos
- npm run build:linux
## 2. Backend
  - cd ..\backend\
  - go mod tidy
  - go build -o ..\frontend\dist\win-unpacked\
    - this changes depending on operating system
## 3. run
  - cd .\dist\win-unpacked\
  - copy pkl files, json files, and py file to this folder 
  - run the python script once before first use
    - py index.py
  - frontend.exe
    - this changes depending on operating system
## 4. to run for dev
  - cd backend
  - go build -o ..\frontend\
  - cd ..\frontend\
  - backend
  - npm run dev
    - this will require a new termnal window