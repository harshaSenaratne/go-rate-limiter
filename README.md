# Rate limiters

- Token Bucket Algorithm
- Per-client rate limiting
- Using tollbooth as middleware

# Demo
- Token Bucket Algorithm
  
![tokenbucket-ezgif com-video-to-gif-converter](https://github.com/user-attachments/assets/5b1acff6-fe55-4e85-bf24-cb82b26fc89a)

- Per-client rate limiting
![perclientrl-ezgif com-video-to-gif-converter](https://github.com/user-attachments/assets/453aebfb-7d96-4441-88db-f1766167e625)

- Using tollbooth as middleware

![tollbooth-ezgif com-video-to-gif-converter](https://github.com/user-attachments/assets/3c61c48e-9245-492f-8ef8-bfb39e7614eb)

## To Run

- cd into the project directory
- run `go run main.go`
- in another terminal, run

## To call the API once - 
```bash
curl -i http://localhost:8080/ping
```
## To call the API multiple times - 

```bash
for i in {1..6}; do curl http://localhost:8080/ping; done
```
